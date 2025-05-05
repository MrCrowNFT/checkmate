package main

import (
	"checkmate/api/internal/handler"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)

	server := &http.Server{
		Addr:    ":8080", //eventually add this to .env
		Handler: mux,
	}

	server.ListenAndServe()
}
