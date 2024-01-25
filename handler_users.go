package main

import (
	"encoding/json"
	"net/http"

	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/models"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := models.UserParams{}

	err := decoder.Decode(&params)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser, err := cfg.DB.CreateUser(params.Email, params.Password)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "cannot create user")
		return
	}

	jsonResponse.ResponedWithJson(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := models.UserParams{}
	err := decoder.Decode(&params)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	loggedUser, err := cfg.DB.UserLogin(params.Email, params.Password)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jsonResponse.ResponedWithJson(w, http.StatusOK, loggedUser)

}
