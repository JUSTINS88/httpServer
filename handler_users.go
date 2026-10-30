package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/JUSTINS88/serverr/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type createUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var dataUser createUserParams
	if err := json.NewDecoder(r.Body).Decode(&dataUser); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(dataUser.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	res, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          dataUser.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Email:     res.Email,
	})
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
