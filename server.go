package main

import (
	"net/http"
)

func start(a args) error {
	f := files{
		dir:             a.dir,
		disallowPersist: a.disallowPersist,
		indexFile:       a.indexFile,
	}
	fileh := handleRequest(&f)
	authh := authorize(a.key, fileh)
	logh := logger(authh)
	s := &http.Server{
		Addr:    a.addr,
		Handler: logh,
	}
	err := s.ListenAndServe()
	return err
}
