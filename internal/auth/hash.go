package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// https://pkg.go.dev/golang.org/x/crypto/bcrypt

// cost should --- 10 <= cost <= 31
// using less for demo purposes
func HashPassword(password string) (string, error) {
	passBuf := []byte(password)
	if len(passBuf) < 3 {
		return "", errors.New("password too short")
	}
	if len(passBuf) > 72 {
		return "", errors.New("password too long")
	}
	retBuf, err := bcrypt.GenerateFromPassword(passBuf, 5)
	if err != nil {
		return "", err
	}
	return string(retBuf), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(password))
}
