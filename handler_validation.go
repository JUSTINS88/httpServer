package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}
	type clean struct {
		CleanedBody string `json:"cleaned_body"`
	}

	defer r.Body.Close()

	var data params
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(data.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	replaceWord(&data.Body, badWords)
	respondWithJSON(w, http.StatusOK, clean{CleanedBody: data.Body})
}

func replaceWord(message *string, badWords map[string]struct{}) {
	msg := strings.Split(*message, " ")
	for i, val := range msg {
		lowerWord := strings.ToLower(val)
		if _, ok := badWords[lowerWord]; ok {
			msg[i] = "****"
		}
	}
	*message = strings.Join(msg, " ")
}
