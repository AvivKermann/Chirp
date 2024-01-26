package jwtauth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(secretKey string, userId int, userExpiresAt *int) (string, error) {
	expiresAt := GetExpiresAt(userExpiresAt)
	claims := jwt.RegisteredClaims{
		Issuer:    fmt.Sprint("chirpy"),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresAt)),
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

func IsTokenValid(token, secretKey string) (*jwt.Token, bool) {

	userToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return &jwt.Token{}, false
	}

	return userToken, true
}

func GetExpiresAt(userExpireInSeconds *int) time.Duration {
	DefaultExpireTimeSeconds := time.Second * 60 * 24

	if userExpireInSeconds == nil {
		return DefaultExpireTimeSeconds
	}
	if *userExpireInSeconds > int(DefaultExpireTimeSeconds) {
		return DefaultExpireTimeSeconds
	}
	if *userExpireInSeconds < 0 {
		return DefaultExpireTimeSeconds
	}

	return time.Duration(*userExpireInSeconds)

}
