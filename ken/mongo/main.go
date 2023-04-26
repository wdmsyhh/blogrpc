package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongotest:27017/?authSource=admin").SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.Database("testdb").CreateCollection(ctx, "testaa"))
}
