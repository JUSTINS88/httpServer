package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/JUSTINS88/serverr/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	var user login
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := cfg.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Hashed PAssword Checking
	pwMatch, err := auth.CheckPasswordHash(user.Password, userData.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !pwMatch {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	//Access Token
	accessToken, err := cfg.makeAccessToken(userData.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Refresh Token
	expiresRefreshToken := time.Now().Add(time.Duration(60*24) * time.Hour)

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userData.ID,
		ExpiresAt: expiresRefreshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal membuat Refresh Token")
		return
	}

	// Response
	respondWithJSON(w, http.StatusOK, loginUser{
		ID:           userData.ID,
		CreatedAt:    userData.CreatedAt,
		UpdatedAt:    userData.UpdatedAt,
		Email:        userData.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) makeAccessToken(userID uuid.UUID) (string, error) {
	durationAccessToken := time.Duration(1) * time.Hour

	accessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, durationAccessToken)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
