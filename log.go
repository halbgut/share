package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		dur := time.Now().Sub(start)
		log.Printf(
			`%v - [%v] "%v" "%v" "%v"`,
			r.RemoteAddr,
			dur,
			fmt.Sprintf("%v %v", r.Method, r.URL.String()),
			r.Referer(),
			r.UserAgent(),
		)
	}
}
