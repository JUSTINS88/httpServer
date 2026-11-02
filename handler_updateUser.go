package main

import (
	"encoding/json"
	"net/http"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/JUSTINS88/serverr/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Validate Access Token
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// json to struct
	defer r.Body.Close()
	var userEmailAndPassword UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&userEmailAndPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Hashing Password
	hashedPassword, err := auth.HashPassword(userEmailAndPassword.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Update Email & Password in DB
	err = cfg.db.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		Email:          userEmailAndPassword.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Gagal mengupdate Email dan Password")
		return
	}

	// Get User information in DB
	userInfo, err := cfg.db.GetUserByEmail(r.Context(), userEmailAndPassword.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email/PAssword Salah")
		return
	}

	respondWithJSON(w, http.StatusOK, UpdateUserInformation{
		ID:        userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email:     userInfo.Email,
	})
}
