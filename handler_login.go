package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/google/uuid"
)

type login struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	var user login
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := cfg.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Hashed PAssword Checking
	hash, err := auth.CheckPasswordHash(user.Password, userData.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hash {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	//JWT
	duration := user.ExpiresInSeconds
	if duration == 0 || duration > 3600 {
		duration = 3600
	}

	JWTUser, err := auth.MakeJWT(userData.ID, cfg.jwtSecret, time.Duration(duration)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	respondWithJSON(w, http.StatusOK, loginUser{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
		Token:     JWTUser,
	})
}
