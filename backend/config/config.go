package config

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    uri := "mongodb://localhost:27017"
    clientOptions := options.Client().ApplyURI(uri)
		log.Printf("Connecting Database: %s", uri)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // Ping the database to verify connection
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    return client, nil
}