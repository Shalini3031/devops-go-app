package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok v1.3"))
	})

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//check img versioning
