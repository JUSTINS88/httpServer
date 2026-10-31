package main

import (
	"net/http"
	"time"

	"github.com/JUSTINS88/serverr/internal/auth"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := cfg.db.CheckRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token salah / tidak ditemukan")
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Silakan login ulang")
		return
	}

	userID := refreshToken.UserID
	accessToken, err := cfg.makeAccessToken(userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type responseToken struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, responseToken{Token: accessToken})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.db.UpdateRevoked(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Gagal mengupdate waktu")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
