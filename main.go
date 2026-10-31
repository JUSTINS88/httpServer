package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/JUSTINS88/serverr/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	Platform := os.Getenv("PLATFORM")
	JWT_Secret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	const path = "."
	const port = "8080"
	cfg := &apiConfig{
		db:        dbQueries,
		platform:  Platform,
		jwtSecret: JWT_Secret,
	}

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
	mux.HandleFunc("POST /admin/reset", cfg.handlerDeleteAllUsers)
	mux.HandleFunc("POST /api/chirps", cfg.createChirps)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeToken)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpsByID)

	serverr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	errs := serverr.ListenAndServe()
	if errs != nil {
		fmt.Println(err)
	}
}
