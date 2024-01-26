package main

import (
	"net/http"

	"github.com/AvivKermann/Chirpy/internal/database"
	"github.com/AvivKermann/Chirpy/internal/jsonResponse"
	"github.com/AvivKermann/Chirpy/internal/jwtauth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token := database.StripPrefix(r.Header.Get("Authorization"))
	isRefreshToken := jwtauth.ValidateRefreshToken(token, cfg.jwtSecret)
	dbExistAndActive := cfg.DB.RefreshTokenExistAndActive(token)
	userId, err := jwtauth.GetIdFromToken(token, cfg.jwtSecret)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if !isRefreshToken || !dbExistAndActive {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := jwtauth.GenerateJwtToken("access", cfg.jwtSecret, userId)
	if err != nil {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "cannot create access token")
		return
	}
	jsonResponse.ResponedWithJson(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken := database.StripPrefix(r.Header.Get("Authorization"))
	revoked := cfg.DB.RevokeRefreshToken(refreshToken)

	if !revoked {
		jsonResponse.ResponedWithError(w, http.StatusBadRequest, "cannot find that token")
	}

	w.WriteHeader(http.StatusOK)
}
