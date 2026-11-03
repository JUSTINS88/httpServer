package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JUSTINS88/serverr/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) hanlderUpdateUserMember(w http.ResponseWriter, r *http.Request) {
	// Check ApiKey
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if polkaKey != cfg.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// declare struct webhook
	type webhookPayload struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// Check identity
	defer r.Body.Close()
	var user webhookPayload
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateChirpyMember(r.Context(), user.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
