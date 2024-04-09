package routes

import (
	"KeDuBak/hashage"
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
		var usersCollection *mongo.Collection

		if c.BodyParser(&dataRequest) != nil && dataRequest.Email != "" &&
			dataRequest.Password != "" && dataRequest.FirstName != "" &&
			dataRequest.LastName != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		usersCollection = client_mongo.Database("kedubak").Collection("User")
		ctx := context.Background()
		filter := bson.M{"email": dataRequest.Email}
		if usersCollection.FindOne(ctx, filter).Decode(&dataUsers) == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Utilisateur déjà existant",
			})
		}
		hash, errHash := hashage.HashPassword(dataRequest.Password)
		if errHash != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		dataRequest.Password = hash
		dataRequest.CreateAt = time.Now()
		dataRequest.LastUpVote = time.Now().Add(-1 * time.Minute)
		user, err := usersCollection.InsertOne(ctx, dataRequest)
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
		var dataRequest structures.User
		var dataUsers structures.User
		var usersCollection *mongo.Collection

		if c.BodyParser(&dataRequest) != nil && dataRequest.Email != "" &&
			dataRequest.Password != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		usersCollection = client_mongo.Database("kedubak").Collection("User")
		ctx := context.Background()
		filter := bson.M{"email": dataRequest.Email}
		if usersCollection.FindOne(ctx, filter).Decode(&dataUsers) != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Utilisateur inexistant",
			})
		}
		if hashage.ComparePasswordWithHash(dataUsers.Password, dataRequest.Password) == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mot de passe incorrect",
			})
		}
		userID := dataUsers.ID.Hex()
		token := jwt_token.GenerateToken(userID)
		if token == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"token": token,
				"user": fiber.Map{
					"email":     dataUsers.Email,
					"firstName": dataUsers.FirstName,
					"lastName":  dataUsers.LastName,
				},
			},
		})
	})
}

func Auth(app *fiber.App, client_mongo *mongo.Client) {
	Register(app, client_mongo)
	Login(app, client_mongo)
}
