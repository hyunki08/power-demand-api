package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// connectDB MongoDB Server
func connectDB() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic("Error loading .env file")
	}

	uri := os.Getenv("MONGO_URI")
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	IsConnected()
}

func disconnectDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

// IsConnected Check MongoDB Client
func IsConnected() bool {
	if client == nil {
		panic("Client is nil")
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return true
}
