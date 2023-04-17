package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rafaeldajuda/database"
	"github.com/rafaeldajuda/entity"
	"github.com/rafaeldajuda/routes"
)

var api entity.API

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.Ctx = context.TODO()
	api.MongoConfig = database.MongoConfig{
		User:       os.Getenv("MONGO_USER"),
		Password:   os.Getenv("MONGO_PASSWORD"),
		Host:       os.Getenv("MONGO_HOST"),
		Port:       os.Getenv("MONGO_PORT"),
		Database:   os.Getenv("MONGO_DATABASE"),
		Collection: os.Getenv("MONGO_COLLECTION"),
	}
	client := api.MongoConfig.CreateConnection(api.Ctx)

	api.HandlerV1.MongoClient = client
	api.HandlerV1.Collection = client.Database(api.MongoConfig.Database).Collection(api.MongoConfig.Collection)
	api.HandlerV1.Ctx = api.Ctx
}

func main() {
	defer func() {
		err := api.MongoClient.Disconnect(api.Ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New()

	// Routes
	routes.Routes(app, &api)

	app.Listen(":8000")

}
