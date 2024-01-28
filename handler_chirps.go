package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/AvivKermann/Chirpy/internal/database"
	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/internal/jwtauth"
	"github.com/go-chi/chi/v5"
)

type parameters struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	token := database.StripPrefix(r.Header.Get("Authorization"))

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
	cleanedContent := getCleanedBody(params.Body)

	userId, err := jwtauth.GetIdFromToken(token, cfg.jwtSecret)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid token")
	}

	newChirp, err := cfg.DB.CreateChirp(cleanedContent, userId)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusInternalServerError, err.Error())
	}
	jsonResponse.ResponedWithJson(w, http.StatusCreated, newChirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chrips, err := cfg.DB.GetChirps()
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, chrips)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	strChirpId := chi.URLParam(r, "chirpId")
	chirpId, err := strconv.Atoi(strChirpId)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, err.Error())
	}
	chirp, exists := cfg.DB.GetSingleChirp(chirpId)

	if !exists {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, "chirp dosen't exist")
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, chirp)

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
