package main

import (
	"net/http"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirps(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

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

	// Get ChripID information in Chirp DB
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found.")
		return
	}

	// Chirp author UUID different with the people who request
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp.")
		return
	}

	if err := cfg.db.DeleteChirps(r.Context(), chirp.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
