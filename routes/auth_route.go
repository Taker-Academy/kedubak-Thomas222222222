package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(app *fiber.App, client_mongo *mongo.Client) {
	app.Post("/auth/register", func(c fiber.Ctx) error {
		fmt.Println("Register !")
		return nil
	})
}

func Login(app *fiber.App, client_mongo *mongo.Client) {
	app.Post("/auth/login", func(c fiber.Ctx) error {
		fmt.Println("Login !")
		return nil
	})
}

func Auth(app *fiber.App, client_mongo *mongo.Client) {
	Register(app, client_mongo)
	Login(app, client_mongo)
}
