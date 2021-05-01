package main

import (
	"net/http"
)

func start(a args) error {
	mux := http.NewServeMux()
	f := newFiles(a.dir)
	mux.HandleFunc("/", handleRequest(&f))
	s := &http.Server{
		Addr:    a.addr,
		Handler: logger(mux),
	}
	err := s.ListenAndServe()
	return err
}
