package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/models"
)

type parameters struct {
	Body string `json:"body"`
}

type response struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "cannot decode chirp")
		return
	}

	valid := validateChirp(params.Body)

	if !valid {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "chirp is invalid")
		return
	}
	newChirp := models.Chirp{}
	cleanedContent := getCleanedBody(params.Body)
	newChirp, err = cfg.DB.CreateChirp(cleanedContent)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusInternalServerError, err.Error())
	}
	jsonResponse.ResponedWithJson(w, http.StatusCreated, response{
		Body: newChirp.Body,
		ID:   newChirp.ID,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chrips, err := cfg.DB.GetChirps()
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, chrips)
}

func getCleanedBody(chirp string) string {

	const replacement = "****"
	words := strings.Split(chirp, " ")
	profane := map[string]struct{}{
		"fornax":    {},
		"kerfuffle": {},
		"sharbert":  {},
	}

	for index, word := range words {
		lowWord := strings.ToLower(word)
		if _, ok := profane[lowWord]; ok {
			words[index] = replacement
		}
	}
	return strings.Join(words, " ")

}

func validateChirp(chirp string) bool {
	if chirpLength := len(chirp); chirpLength > 140 {
		return false
	}
	return true

}
