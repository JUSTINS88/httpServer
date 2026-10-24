package main

import (
	"fmt"
	"net/http"
)

func main() {
	const path = "."
	const port = "8080"
	cfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/healthz",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	mux.Handle("/app/", cfg.middlewareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir(path)))))
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidation)

	serverr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err := serverr.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
