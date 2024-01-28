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
	Accesstoken, err := jwtauth.GenerateJwtToken("access", cfg.jwtSecret, loggedUser.ID)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "cannot create access token")
		return
	}
	refreshToken, err := jwtauth.GenerateJwtToken("refresh", cfg.jwtSecret, loggedUser.ID)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "cannot create refresh token")
		return
	}
	err = cfg.DB.CreateNewRefreshToken(refreshToken)

	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusInternalServerError, "cannot insert refresh token")
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, struct {
		models.ResponseUser
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ResponseUser: loggedUser,
		Token:        Accesstoken,
		RefreshToken: refreshToken,
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

	tokenType, err := jwtauth.GetIssuerFromToken(token, cfg.jwtSecret)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	if tokenType == "chirpy-refresh" {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid action for refresh token")
		return
	}

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
func (cfg *apiConfig) handlerUserSubscribe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := models.WebhooksParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if event := params.Event; event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}
	user, err := cfg.DB.GetUserById(params.Data.UserID)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusNotFound, "user not found")
		return
	}

	upgraded := cfg.DB.MakeUserChirpyRed(user.ID)
	if !upgraded {
		jsonResponse.ResponedWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, struct {
		Body string `json:"body"`
	}{
		Body: "",
	})
}
