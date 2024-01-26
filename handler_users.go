package main

import (
	"encoding/json"
	"net/http"

	"github.com/AvivKermann/Chirpy/internal/database"
	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/internal/jwtauth"
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
	token, err := jwtauth.GenerateJwtToken(cfg.jwtSecret, loggedUser.ID, params.ExpireInSeconds)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "cannot create token")
		return
	}

	jsonResponse.ResponedWithJson(w, http.StatusOK, struct {
		models.ResponseUser
		Token string `json:"token"`
	}{
		ResponseUser: loggedUser,
		Token:        token,
	})

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := models.UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	token := database.StripPrefix(r.Header.Get("Authorization"))

	userId, err := jwtauth.GetIdFromToken(token, cfg.jwtSecret)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := cfg.DB.GetUserById(userId)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusInternalServerError, "cannot find user")
		return
	}

	updatedUser, err := cfg.DB.UpdateUser(params.Email, params.Password, user)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "cannot change user")
		return
	}

	jsonResponse.ResponedWithJson(w, http.StatusOK, updatedUser)

}
