package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27012/?authSource=admin").SetAuth(options.Credential{
	//	Username: "root",
	//	Password: "root",
	//}))
	// 连接容器副本集时加上 connect=direct
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27016/?connect=direct"))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(client.Database("testdb").CreateCollection(ctx, "cola"))

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	col := client.Database("testdb").Collection("cola")
	id, _ := primitive.ObjectIDFromHex("6480319cfbcabc003bd46004")
	selector := primitive.M{
		"_id":       id,
		"isDeleted": false,
	}
	updater := primitive.M{
		"$set": primitive.M{
			"isDeleted": true,
		},
	}
	one, err := col.UpdateOne(ctx, selector, updater)
	log.Println("one", one)
	log.Println("err", err)
}
