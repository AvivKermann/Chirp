package jwtauth

import (
	"fmt"
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
