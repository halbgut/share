package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
)

type files struct {
	dir             string
	waiters         sync.Map
	disallowPersist bool
}

type waiter struct {
	error chan error
	data  chan []byte
}

func (f *files) Post(ctx context.Context, rpath string, r io.ReadCloser, persist bool) error {
	defer r.Close()
	if err := validate(rpath); err != nil {
		return fmt.Errorf("Invalid path: %w", err)
	}
	if !f.disallowPersist && persist {
		err := f.postToFile(ctx, rpath, r)
		return err
	}
	buf := make([]byte, 4*1024)
	n, err := r.Read(buf)
	if errors.Is(err, io.EOF) {
		return f.postChunkToFile(ctx, rpath, buf[:n])
	}
	if err != nil {
		return fmt.Errorf("Failed to read request: %w", err)
	}
	return f.postToWaiter(ctx, rpath, buf, n, r)
}

func (f *files) postChunkToFile(ctx context.Context, rpath string, buf []byte) error {
	fp := path.Join(f.dir, rpath)
	err := os.WriteFile(fp, buf, 0666)
	if err != nil {
		return fmt.Errorf("Failed to write file: %w", err)
	}
	return nil
}

func (f *files) postToFile(ctx context.Context, rpath string, r io.Reader) error {
	fp := path.Join(f.dir, rpath)
	fl, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("Failed to create new file: %w", err)
	}
	_, err = io.Copy(fl, r)
	if err != nil {
		return fmt.Errorf("Failed to write file: %w", err)
	}
	return nil
}

func (f *files) postToWaiter(ctx context.Context, rpath string, buf []byte, end int, r io.Reader) error {
	total := len(buf[:end])
	_, ok := f.waiters.Load(rpath)
	if ok {
		return fmt.Errorf("Conflict. Path already currently in use")
	}
	var wait waiter
	wait.error = make(chan error)
	wait.data = make(chan []byte)
	defer close(wait.error)
	defer close(wait.data)
	defer f.waiters.Delete(rpath)
	f.waiters.Store(rpath, wait)
	wait.data <- buf[:end]
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		n, err := r.Read(buf)
		wait.data <- buf[:n]
		total += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			err := fmt.Errorf("Failed to read body: %w", err)
			wait.error <- err
			return err
		}
	}
}

func (f *files) Get(ctx context.Context, rpath string, w io.Writer) error {
	if err := validate(rpath); err != nil {
		return fmt.Errorf("Invalid path: %w", err)
	}
	iw, ok := f.waiters.Load(rpath)
	if ok {
		wait := iw.(waiter)
		return f.getFromWaiter(ctx, wait, w)
	} else {
		return f.getFromFile(ctx, rpath, w)
	}
}

func (f *files) getFromFile(ctx context.Context, rpath string, w io.Writer) error {
	fp := path.Join(f.dir, rpath)
	fl, err := os.Open(fp)
	if _, ok := err.(*fs.PathError); ok {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	_, err = io.Copy(w, fl)
	if err != nil {
		return fmt.Errorf("Failed to send response: %w", err)
	}
	return nil
}

func (f *files) getFromWaiter(ctx context.Context, wait waiter, w io.Writer) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		b, ok := <-wait.data
		_, err := w.Write(b)
		if err != nil {
			return fmt.Errorf("Failed to write response: %w", err)
		}
		if !ok {
			break
		}
	}
	err := <-wait.error
	return err
}
