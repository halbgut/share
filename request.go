package main

import (
	"errors"
	"fmt"
	"net/http"
)

func handleRequest(f *files) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var err error
		path := r.URL.Path
		switch r.Method {
		case http.MethodGet:
			err = f.Get(ctx, path, w)
		case http.MethodPost:
			err = f.Post(ctx, path, r.Body)
			if err == nil {
				_, err = w.Write([]byte{})
			}
		default:
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := fmt.Sprintf("Request failed: %v\n", err)
			w.Write([]byte(res))
		}
	}
}
