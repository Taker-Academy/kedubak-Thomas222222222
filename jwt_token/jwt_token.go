package jwt_token

import (
	"KeDuBak/structures"

	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckToken(tokenString string, client_mongo *mongo.Client) (string, int) {
	var userID string
	var dataUsers structures.User

	secretKey := []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return userID, -1
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["id"].(string)
		usersCollection := client_mongo.Database("kedubak").Collection("User")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return userID, -1
		}
		ctx := context.Background()
		filter := bson.M{"_id": objectID}
		if err := usersCollection.FindOne(ctx, filter).Decode(&dataUsers); err != nil {
			return userID, -1
		}
		if userID != "" {
			return userID, 0
		}
	}
	return userID, -1
}

func GenerateToken(userID string) string {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	secretKey := []byte(os.Getenv("SECRET"))
	FirstToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	FinalToken, err := FirstToken.SignedString(secretKey)
	if err != nil {
		return ""
	}
	return FinalToken
}
