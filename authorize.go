package main

import (
	"net/http"
)

func authorize(key string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		q := r.URL.Query()
		var k string
		if ks := q["key"]; len(ks) > 0 {
			k = ks[0]
		}
		if k != key {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access denied!"))
			return
		}
		next.ServeHTTP(w, r)
	}
}
