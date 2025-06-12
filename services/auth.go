package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword membandingkan hash dan plain password
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWT membuat JWT signed dengan HS256
func GenerateJWT(userID uint, email string) (string, error) {
	// secret := os.Getenv("JWT_SECRET", "azsxdcfv")
	secret := "azsxdcfv"
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	exp := time.Now().Add(1 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
