package jwtauth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(tokenType, secretKey string, userId int) (string, error) {
	const RefreshDefaultExpireTimeDay = 60
	const AccessDefaultExpireTimeHour = 1
	expiresAt := AccessDefaultExpireTimeHour

	if tokenType == "refresh" {
		expiresAt = RefreshDefaultExpireTimeDay * 24
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-" + tokenType,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * time.Duration(expiresAt))),
		Subject:   fmt.Sprint(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
func GetIdFromToken(token, secretKey string) (int, error) {

	userToken, valid := IsTokenValid(token, secretKey)
	if !valid {
		return 0, errors.New("invalid token")
	}
	stringId, err := userToken.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	intId, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, err
	}
	return intId, nil

}

func GetIssuerFromToken(token, secretKey string) (string, error) {
	userToken, valid := IsTokenValid(token, secretKey)
	if !valid {
		return "", errors.New("invalid token")
	}
	issuer, err := userToken.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	return issuer, nil
}

func IsTokenValid(token, secretKey string) (*jwt.Token, bool) {

	userToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return &jwt.Token{}, false
	}

	return userToken, true
}

func ValidateRefreshToken(refreshToken, secretKey string) bool {
	token, isValidToken := IsTokenValid(refreshToken, secretKey)
	if !isValidToken {
		return false
	}

	if tokenType, err := token.Claims.GetIssuer(); tokenType != "chirpy-refresh" || err != nil {
		return false
	}

	return true

}
