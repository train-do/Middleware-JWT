package utils

import (
	"errors"
	"time"
	"voucher_system/config"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey []byte

// InitJwtKey membaca JWT_KEY dari environment sekali saja
func InitJwtKey(config config.Configuration) error {
	key := config.JwtKey
	if key == "" {
		return errors.New("JWT_KEY is not set in the environment")
	}
	JwtKey = []byte(key)
	return nil
}

func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   string(rune(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}