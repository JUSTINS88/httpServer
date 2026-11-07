package main

import (
	"net/http"
	"slices"

	"github.com/JUSTINS88/serverr/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")

	var res []database.Chirp
	var err error
	if author != "" {
		authorUUID, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err = cfg.db.GetAllChirpsByAuthorID(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Cannot load Chrips")
			return
		}
	} else {
		res, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Cannot load Chrips")
			return
		}
	}

	resp := make([]Chirp, 0, len(res))
	for _, dbChrip := range res {
		resp = append(resp, Chirp{
			ID:        dbChrip.ID,
			CreatedAt: dbChrip.CreatedAt,
			UpdatedAt: dbChrip.UpdatedAt,
			Body:      dbChrip.Body,
			UserID:    dbChrip.UserID,
		})
	}

	// Sorting Chirp
	sortChirp := r.URL.Query().Get("sort")
	if sortChirp == "desc" {
		slices.SortFunc(resp, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	} else {
		slices.SortFunc(resp, func(a, b Chirp) int {
			return a.CreatedAt.Compare(b.CreatedAt)
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
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        msg.ID,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Body:      msg.Body,
		UserID:    msg.UserID,
	})
}
