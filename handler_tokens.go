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

	if !isRefreshToken || !dbExistAndActive {
		jsonResponse.ResponedWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	jsonResponse.ResponedWithJson(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
