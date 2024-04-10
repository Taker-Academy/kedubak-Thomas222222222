package routes

import (
	"KeDuBak/hashage"
	"KeDuBak/jwt_token"
	"KeDuBak/structures"
	"fmt"

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

func CheckValidityOfInfo(c *fiber.Ctx, dataUser *structures.User, usersCollection *mongo.Collection) int {
	var dataRequest structures.User

	if c.BodyParser(&dataRequest) != nil || dataRequest.Email == "" ||
		dataRequest.Password == "" {
		return -1
	}
	hash, errHash := hashage.HashPassword(dataRequest.Password)
	if errHash != nil {
		return -1
	}
	dataUser.Email = dataRequest.Email
	dataUser.FirstName = dataRequest.FirstName
	dataUser.LastName = dataRequest.LastName
	dataUser.Password = hash
	ctx := context.Background()
	filter := bson.M{"email": dataRequest.Email}
	err := usersCollection.FindOne(ctx, filter).Decode(&dataRequest)
	if err == nil && dataUser.ID != dataRequest.ID {
		return -1
	}
	return 0
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
		usersCollection := client_mongo.Database("kedubak").Collection("User")
		if CheckValidityOfInfo(c, &dataUser, usersCollection) == -1 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"ok":    false,
				"error": "Echec de validation des paramètres",
			})
		}
		objectID, errID := primitive.ObjectIDFromHex(userID)
		if errID != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		ctx := context.Background()
		filter := bson.M{"_id": objectID}
		update := bson.M{
			"$set": bson.M{
				"email":     dataUser.Email,
				"password":  dataUser.Password,
				"firstName": dataUser.FirstName,
				"lastName":  dataUser.LastName,
			},
		}
		if _, err := usersCollection.UpdateOne(ctx, filter, update); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"ok":    false,
				"error": "Echec de validation des paramètres",
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

func DeletePost(userID string, client_mongo *mongo.Client) int {
	postCollection := client_mongo.Database("kedubak").Collection("Post")
	ctx := context.Background()
	cursor, err := postCollection.Find(ctx, bson.M{})
	if err != nil {
		return -1
	}
	for cursor.Next(ctx) {
		var post structures.Post
		if err := cursor.Decode(&post); err != nil {
			return -1
		}
		filter := bson.M{"_id": post.ID}
		if string(post.UserID) == userID {
			if _, errDelete := postCollection.DeleteOne(ctx, filter); errDelete != nil {
				return -1
			}
		}
	}
	return 0
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
		if DeletePost(userID, client_mongo) == -1 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Utilisateur non trouvé",
			})
		}
		if _, errDelete := usersCollection.DeleteOne(ctx, filter); errDelete != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Utilisateur non trouvé",
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
