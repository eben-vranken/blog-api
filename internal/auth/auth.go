package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	return string(bytes), err
}

func CheckPassword(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))

	return err == nil
}
