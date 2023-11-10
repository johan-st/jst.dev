package main

import (
	"net/http"
	"time"
)

func main() {
	l := newLogger()


	handler := routes(l)

	server := http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
		ErrorLog:          l,
	}

	l.Println("Starting server on port 8080")
	server.ListenAndServe()
}
