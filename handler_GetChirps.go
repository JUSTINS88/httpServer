package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp []Chrip
	for i, _ := range res {
		resp = append(resp, Chrip{
			ID:        res[i].ID,
			CreatedAt: res[i].CreatedAt,
			UpdatedAt: res[i].UpdatedAt,
			Body:      res[i].Body,
			UserID:    res[i].UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) GetChirpsByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := cfg.db.GetChirp(r.Context(), uuid)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, Chrip{
		ID:        msg.ID,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Body:      msg.Body,
		UserID:    msg.UserID,
	})
}
