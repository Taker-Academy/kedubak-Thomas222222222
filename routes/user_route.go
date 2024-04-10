package routes

import (
	"KeDuBak/jwt_token"
	"KeDuBak/structures"

	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDataUser(dataUser *structures.User, client_mongo *mongo.Client, userID string) int {
	usersCollection := client_mongo.Database("kedubak").Collection("User")
	objectID, errID := primitive.ObjectIDFromHex(userID)
	if errID != nil {
		return -1
	}
	ctx := context.Background()
	filter := bson.M{"_id": objectID}
	if err := usersCollection.FindOne(ctx, filter).Decode(dataUser); err != nil || dataUser.Email == "" ||
		dataUser.FirstName == "" || dataUser.LastName == "" {
		return -1
	}
	return 0
}

func Me(app *fiber.App, client_mongo *mongo.Client) {
	app.Get("/user/me", func(c *fiber.Ctx) error {
		var dataUser structures.User
		var userID string
		var errToken int

		token := c.Get("Authorization")
		if userID, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		if GetDataUser(&dataUser, client_mongo, userID) == -1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"email":     dataUser.Email,
				"firstName": dataUser.FirstName,
				"lastName":  dataUser.LastName,
			},
		})
	})
}

func Edit(app *fiber.App, client_mongo *mongo.Client) {
	app.Put("/user/edit", func(c *fiber.Ctx) error {
		var dataUser structures.User
		var userID string
		var errToken int

		token := c.Get("Authorization")
		if userID, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		if GetDataUser(&dataUser, client_mongo, userID) == -1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"email":     dataUser.Email,
				"firstName": dataUser.FirstName,
				"lastName":  dataUser.LastName,
			},
		})
	})
}

func Delete(app *fiber.App, client_mongo *mongo.Client) {
	app.Delete("/user/remove", func(c *fiber.Ctx) error {
		var dataUser structures.User
		var userID string
		var errToken int

		token := c.Get("Authorization")
		if userID, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		if GetDataUser(&dataUser, client_mongo, userID) == -1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		usersCollection := client_mongo.Database("kedubak").Collection("User")
		objectID, errID := primitive.ObjectIDFromHex(userID)
		if errID != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		ctx := context.Background()
		filter := bson.M{"_id": objectID}
		if _, errDelete := usersCollection.DeleteOne(ctx, filter); errDelete != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":		false,
				"error":	"Utilisateur non trouvé",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"email":     dataUser.Email,
				"firstName": dataUser.FirstName,
				"lastName":  dataUser.LastName,
				"remove":    true,
			},
		})
	})
}

func User(app *fiber.App, client_mongo *mongo.Client) {
	Me(app, client_mongo)
	Edit(app, client_mongo)
	Delete(app, client_mongo)
}
