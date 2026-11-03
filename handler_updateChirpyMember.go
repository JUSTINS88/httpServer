package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type webhookPayload struct {
	Event string   `json:"event"`
	Data  userData `json:"data"`
}
type userData struct {
	UserID string `json:"user_id"`
}

func (cfg *apiConfig) hanlderUpdateUserMember(w http.ResponseWriter, r *http.Request) {
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

	userID, err := uuid.Parse(user.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = cfg.db.UpdateChirpyMember(r.Context(), userID)
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
