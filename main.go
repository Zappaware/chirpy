package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static assets
	assets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// Serve index.html (optional, if you already have this)
	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

