package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "./app"
	const port = "8080"
	fs := http.FileServer(http.Dir("./app"))

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/app/", fs))
	mux.HandleFunc("/healthz/", readinessHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	const body = "OK"
	const contentTypeKey = "Content-Type"
	const contentTypeValue = "text/plain; charset=utf-8"
	const statusCode = 200
	const resBody = "OK"

	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(resBody))
	if err != nil {
		fmt.Printf("error writing response: %v", err)
	}
}
