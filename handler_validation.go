package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}
	type validValue struct {
		Valid bool `json:"valid"`
	}

	defer r.Body.Close()

	var data params
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(data.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp is too long"))
		return
	}

	respondWithJSON(w, http.StatusOK, validValue{Valid: true})
}
