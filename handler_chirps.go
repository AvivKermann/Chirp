package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/AvivKermann/Chirpy/internal/database"
	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/internal/jwtauth"
	"github.com/AvivKermann/Chirpy/models"
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
	strAuthorId := r.URL.Query().Get("author_id")
	defaultOrderType := "asc"
	orderType := r.URL.Query().Get("sort")

	if orderType != "desc" {
		orderType = defaultOrderType
	}

	if strAuthorId != "" {
		authorId, err := strconv.Atoi(strAuthorId)
		if err != nil {
			jsonResponse.ResponedWithError(w, http.StatusBadRequest, "invalid id")
			return
		}
		chirps, err := cfg.DB.GetChirpsByAuthor(authorId)
		if err != nil {
			jsonResponse.ResponedWithError(w, http.StatusBadRequest, "invalid id")
			return
		}
		SortChirps(orderType, chirps)
		jsonResponse.ResponedWithJson(w, http.StatusOK, chirps)
		return
	}

	chrips, err := cfg.DB.GetChirps()
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, err.Error())
		return
	}
	SortChirps(orderType, chrips)
	jsonResponse.ResponedWithJson(w, http.StatusOK, chrips)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	strChirpId := chi.URLParam(r, "chirpId")
	chirpId, err := strconv.Atoi(strChirpId)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, exists := cfg.DB.GetSingleChirp(chirpId)

	if !exists {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, "chirp dosen't exist")
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	strChirpId := chi.URLParam(r, "chirpId")
	chirpId, err := strconv.Atoi(strChirpId)
	token := database.StripPrefix(r.Header.Get("Authorization"))
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := jwtauth.GetIdFromToken(token, cfg.jwtSecret)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "not a user")
		return
	}

	chirp, exist := cfg.DB.GetSingleChirp(chirpId)
	if !exist {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "chirp dosent exist")
		return
	}

	isOwner := isChirpAuthor(chirp, userId)
	if !isOwner {
		jsonResponse.ResponedWithError(w, http.StatusForbidden, "cannot delete chirp by other users")
		return
	}

	isDeleted := cfg.DB.DeleteSingleChirp(chirpId)
	if !isDeleted {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "chirp cannot be deleted")
		return
	}

	w.WriteHeader(http.StatusOK)
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

func isChirpAuthor(chirp models.Chirp, userId int) bool {
	return userId == chirp.AuthorID
}

func SortChirps(sortingType string, chirps []models.Chirp) {
	if sortingType == "asc" {
		sortChirpAsc(chirps)
	} else if sortingType == "desc" {
		sortChirpDesc(chirps)
	}
}

func sortChirpAsc(chirps []models.Chirp) {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID

	})
}
func sortChirpDesc(chirps []models.Chirp) {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID > chirps[j].ID
	})
}
