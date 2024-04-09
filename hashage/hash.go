package hashage

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswordWithHash(hash string, password string) int {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return -1
	}
	return 0
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}
