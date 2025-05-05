package main

import (
	"checkmate/api/internal/handler"
	"checkmate/api/internal/storage"
	"net/http"
)

func main() {

	storage.InitDb()

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)

	server := &http.Server{
		Addr:    ":8080", //eventually add this to .env
		Handler: mux,
	}

	server.ListenAndServe()
}
