package services

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword membandingkan hash dan plain password
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWT membuat JWT signed dengan HS256
func GenerateJWT(userID uint, email string) (string, int64, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", 0, fmt.Errorf("JWT_SECRET not set")
	}

	exp := time.Now().Add(1 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, exp, nil
}
