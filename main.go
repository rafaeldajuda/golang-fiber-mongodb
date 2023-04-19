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

	api.Name = os.Getenv("PROJECT_NAME")
	api.Host = os.Getenv("PROJECT_HOST")
	api.Port = os.Getenv("PROJECT_PORT")

	api.Ctx = context.TODO()
	api.MongoConfig = database.MongoConfig{
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
	client := api.MongoConfig.CreateConnection(api.Ctx)

	api.HandlerV1.Ctx = api.Ctx
	api.HandlerV1.Database = os.Getenv("MONGO_DATABASE")
	api.HandlerV1.Collection = os.Getenv("MONGO_COLLECTION")
	api.HandlerV1.MongoClient = client
	api.HandlerV1.MongoCollection = client.Database(api.HandlerV1.Database).Collection(api.HandlerV1.Collection)
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

	err := app.Listen(api.Host + ":" + api.Port)
	if err != nil {
		log.Fatal(err)
	}
}
