package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func error_hanling(MongoURL *string, SECRET *string) int {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error : file .env not found\n")
		return -1
	}
	*MongoURL = os.Getenv("MONGO_URL")
	*SECRET = os.Getenv("SECRET")
	if *SECRET == "" || *MongoURL == "" {
		fmt.Print("Error : missing environment variables\n")
		return -1
	}
	return 0
}

func main() {
	var MongoURL string
	var SECRET string

	if error_hanling(&MongoURL, &SECRET) == -1 {
		os.Exit(1)
	}
	fmt.Printf("Mongo URL : %s\n", MongoURL)
	fmt.Printf("Secret : %s\n", SECRET)
}
