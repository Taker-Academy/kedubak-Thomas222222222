package routes

import (
	"KeDuBak/hashages"
	"KeDuBak/jwt_token"
	"KeDuBak/structures"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(app *fiber.App, client_mongo *mongo.Client) {
	app.Post("/auth/register", func(c *fiber.Ctx) error {
		var dataRequest structures.User
		var dataUsers structures.User
		var userCollection *mongo.Collection

		if c.BodyParser(&dataRequest) != nil && dataRequest.Email != "" &&
			dataRequest.Password != "" && dataRequest.FirstName != "" &&
			dataRequest.LastName != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		userCollection = client_mongo.Database("kedubak").Collection("User")
		ctx := context.Background()
		filter := bson.M{"email": dataRequest.Email}
		if userCollection.FindOne(ctx, filter).Decode(&dataUsers) == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Utilisateur déjà existant",
			})
		}
		hash, errHash := hashages.HashPassword(dataRequest.Password)
		if errHash != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		dataRequest.Password = hash
		dataRequest.CreateAt = time.Now()
		dataRequest.LastUpVote = time.Now().Add(-1 * time.Minute)
		user, err := userCollection.InsertOne(ctx, dataRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		userID := user.InsertedID.(primitive.ObjectID).Hex()
		token := jwt_token.GenerateToken(userID)
		if token == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"token": token,
				"user": fiber.Map{
					"email":     dataRequest.Email,
					"firstName": dataRequest.FirstName,
					"lastName":  dataRequest.LastName,
				},
			},
		})
	})
}

func Login(app *fiber.App, client_mongo *mongo.Client) {
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": "Erreur interne du serveur",
		})
	})
}

func Auth(app *fiber.App, client_mongo *mongo.Client) {
	Register(app, client_mongo)
	Login(app, client_mongo)
}
