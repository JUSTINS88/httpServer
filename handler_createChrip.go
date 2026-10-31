package main

import (
	"encoding/json"
	"net/http"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/JUSTINS88/serverr/internal/database"
)

func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	var sender Sender
	if err := json.NewDecoder(r.Body).Decode(&sender); err != nil {
		respondWithError(w, http.StatusBadRequest, "Gagal memparse JSON")
		return
	}

	// Checking JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Take data from DB
	res, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   sender.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal mengambil dari database")
		return
	}

	// Response
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	})
}
