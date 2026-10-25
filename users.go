package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Reqq struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var data Reqq
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := cfg.db.CreateUser(r.Context(), data.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resJson := Users{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Email:     res.Email,
	}

	respondWithJSON(w, http.StatusCreated, resJson)
}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Error bang")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type success struct {
		Success bool `json:"success"`
	}
	respondWithJSON(w, http.StatusOK, success{Success: true})
}
