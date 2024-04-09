package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"KeDuBak/database"
	"KeDuBak/routes"
)

func error_hanling(MongoURL *string) int {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error : file .env not found\n")
		return -1
	}
	*MongoURL = os.Getenv("MONGO_URL")
	if *MongoURL == "" {
		fmt.Print("Error : missing environment variables\n")
		return -1
	}
	return 0
}

func main() {
	var MongoURL string

	if error_hanling(&MongoURL) == -1 {
		os.Exit(1)
	}
	app := fiber.New()
	client_mongo := database.ConnectDB(MongoURL)
	defer func() {
		if err := client_mongo.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	app.Use(cors.New())
	routes.Auth(app, client_mongo)
	routes.Me(app, client_mongo)
	app.Listen(":8080")
}
