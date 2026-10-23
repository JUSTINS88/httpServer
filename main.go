package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func main() {
	const path = "."
	const port = "8080"
	cfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /healthz",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

	mux.Handle("/app/", cfg.middlewareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir(path)))))
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /reset", cfg.resetMetrics)

	serverr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err := serverr.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
