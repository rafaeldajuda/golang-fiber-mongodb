package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	Collection string
}

func (mc *MongoConfig) CreateConnection(ctx context.Context) *mongo.Client {
	// mongodb uri format
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", mc.User, mc.Password, mc.Host, mc.Port)
	if uri == "" {
		log.Fatal("You must set your 'uri' variable.")
	}

	// connection
	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("MongoDB connection error: " + err.Error())
	}

	// check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
