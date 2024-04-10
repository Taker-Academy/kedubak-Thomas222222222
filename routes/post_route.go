package routes

import (
	"KeDuBak/jwt_token"
	"KeDuBak/structures"

	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DisplayMe(app *fiber.App, client_mongo *mongo.Client) {
	app.Get("/post/me", func(c *fiber.Ctx) error {
		var listPosts []structures.Post
		var userID string
		var errToken int

		token := c.Get("Authorization")
		if userID, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		postCollection := client_mongo.Database("kedubak").Collection("Post")
		ctx := context.Background()
		cursor, err := postCollection.Find(ctx, bson.M{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		for cursor.Next(ctx) {
			var post structures.Post
			if err := cursor.Decode(&post); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ok":    false,
					"error": "Erreur interne du serveur",
				})
			}
			if string(post.UserID) == userID {
				listPosts = append(listPosts, post)
			}
		}
		if listPosts == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Aucun post trouver pour cette utilisateur",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"data": listPosts,
		})
	})
}

func DetailsPost(app *fiber.App, client_mongo *mongo.Client) {
	app.Get("/post/:id", func(c *fiber.Ctx) error {
		var post structures.Post
		var errToken int

		token := c.Get("Authorization")
		if _, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		objectID, errObjectID := primitive.ObjectIDFromHex(c.Params("id"))
		if errObjectID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		postCollection := client_mongo.Database("kedubak").Collection("Post")
		ctx := context.Background()
		filter := bson.M{"_id": objectID}
		err := postCollection.FindOne(ctx, filter).Decode(&post)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Élément non trouvé",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"data": post,
		})
	})
}

func DeleteSpecificPost(app *fiber.App, client_mongo *mongo.Client) {
	app.Delete("/post/:id", func(c *fiber.Ctx) error {
		var userID string
		var post structures.Post
		var errToken int

		token := c.Get("Authorization")
		if userID, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		objectID, errObjectID := primitive.ObjectIDFromHex(c.Params("id"))
		if errObjectID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		postCollection := client_mongo.Database("kedubak").Collection("Post")
		ctx := context.Background()
		filter := bson.M{"_id": objectID}
		err := postCollection.FindOne(ctx, filter).Decode(&post)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Élément non trouvé",
			})
		}
		if string(post.UserID) != userID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "L'utilisateur n'est pas le propriétaire de l'élément",
			})
		}
		if _, errDelete := postCollection.DeleteOne(ctx, filter); errDelete != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Élément non trouvé",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"_id":       post.ID,
				"createdAt": post.CreateAt,
				"userId":    post.UserID,
				"firstName": post.FirstName,
				"title":     post.Title,
				"content":   post.Content,
				"comments":  post.Comments,
				"upVotes":   post.UpVotes,
				"removed":   true,
			},
		})
	})
}

func Display(app *fiber.App, client_mongo *mongo.Client) {
	app.Get("/post", func(c *fiber.Ctx) error {
		var listPosts []structures.Post
		var errToken int

		token := c.Get("Authorization")
		if _, errToken = jwt_token.CheckToken(token, client_mongo); errToken == -1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvais token JWT",
			})
		}
		postCollection := client_mongo.Database("kedubak").Collection("Post")
		ctx := context.Background()
		cursor, err := postCollection.Find(ctx, bson.M{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		for cursor.Next(ctx) {
			var post structures.Post
			if err := cursor.Decode(&post); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ok":    false,
					"error": "Erreur interne du serveur",
				})
			}
			listPosts = append(listPosts, post)
		}
		if listPosts == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"data": listPosts,
		})
	})
}

func Create(app *fiber.App, client_mongo *mongo.Client) {
	app.Post("/post", func(c *fiber.Ctx) error {
		var postRequest structures.Post
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
		if c.BodyParser(&postRequest) != nil || postRequest.Content == "" ||
			postRequest.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		if GetDataUser(&dataUser, client_mongo, userID) == -1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Mauvaise requête, paramètres manquants ou invalides",
			})
		}
		postCollection := client_mongo.Database("kedubak").Collection("Post")
		newPost := structures.Post{
			CreateAt:  time.Now(),
			UserID:    userID,
			FirstName: dataUser.FirstName,
			Title:     postRequest.Title,
			Content:   postRequest.Content,
			Comments:  []structures.Comments{},
			UpVotes:   []string{},
		}
		ctx := context.Background()
		_, err := postCollection.InsertOne(ctx, newPost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Erreur interne du serveur",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"ok":   true,
			"data": newPost,
		})
	})
}

func Post(app *fiber.App, client_mongo *mongo.Client) {
	DisplayMe(app, client_mongo)
	DetailsPost(app, client_mongo)
	DeleteSpecificPost(app, client_mongo)
	Display(app, client_mongo)
	Create(app, client_mongo)
}
