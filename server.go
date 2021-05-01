package main

import (
	"net/http"
)

func start(a args) error {
	mux := http.NewServeMux()
	f := files{
		dir:             a.dir,
		disallowPersist: a.disallowPersist,
	}
	mux.HandleFunc("/", handleRequest(&f))
	s := &http.Server{
		Addr:    a.addr,
		Handler: logger(mux),
	}
	err := s.ListenAndServe()
	return err
}
