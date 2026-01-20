package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	if _, err := fmt.Fprintf(w, "Hits: %d", hits); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", readinessHandler)

	fs := http.FileServer(http.Dir("."))
	handler := http.StripPrefix("/app", fs)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	// Serve static assets
	assets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
