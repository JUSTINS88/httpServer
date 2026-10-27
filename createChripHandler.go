package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JUSTINS88/serverr/internal/database"
	"github.com/google/uuid"
)

type Sender struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chrips struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChrips(w http.ResponseWriter, r *http.Request) {
	var sender Sender
	if err := json.NewDecoder(r.Body).Decode(&sender); err != nil {
		respondWithError(w, http.StatusBadRequest, "Gagal memparse JSON")
		return
	}

	res, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   sender.Body,
		UserID: sender.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal mengambil dari database")
		return
	}

	resp := Chrips{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}
