package jwt_token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) string {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	secretKey := []byte(os.Getenv("SECRET"))
	Token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	finalToken, err := Token1.SignedString(secretKey)
	if err != nil {
		return ""
	}
	return finalToken
}
