package main

import (
	"log"
	"net/http"
	"strings"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Host, "localhost") {
			log.Println("Log request for route:", r.RequestURI)
		} else {
			go logRequestToPlausible(NewLogRequestBody(r.RequestURI, r.Header.Get("X-Forwarded-For")))
		}
		next.ServeHTTP(w, r)
	})
}
