package main

import (
	"encoding/json"
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
	html := fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func respondWithError(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(data))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(data))
}

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
	mux.HandleFunc(
		"POST /api/validate_chirp",
		func(w http.ResponseWriter, r *http.Request) {
			type params struct {
				Body string `json:"body"`
			}
			type valids struct {
				Valid bool `json:"valid"`
			}
			type errorJson struct {
				Error string `json:"error"`
			}

			defer r.Body.Close()

			var data params
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			if len(data.Body) > 140 {
				var eror errorJson
				eror.Error = "Chirp is too long"
				respondWithError(w, http.StatusBadRequest, eror)
				return
			}

			var validd valids
			validd.Valid = true
			respondWithJSON(w, http.StatusOK, validd)
			return
		})

	serverr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err := serverr.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
