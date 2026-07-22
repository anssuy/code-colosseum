package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)

	return err == nil
}
