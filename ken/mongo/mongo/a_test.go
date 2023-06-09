package mongo

import (
	"blogrpc/core/extension/bson"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (self *PlaygroundSuite) TestGetOne() {
	ctx := context.Background()
	col := client.Database("testdb").Collection("cola")
	id, _ := primitive.ObjectIDFromHex("6480319cfbcabc003bd46004")
	selector := primitive.M{
		"_id": id,
	}
	one := col.FindOne(ctx, selector)
	result := bson.M{}
	one.Decode(&result)
	log.Println("===============")
	log.Println(result)
}

func (self *PlaygroundSuite) TestGetOnev2() {
	ctx := context.Background()
	col := client.Database("testdb").Collection("cola")
	id := bson.ObjectIdHex("6480319cfbcabc003bd46004")
	selector := bson.M{
		"_id": id,
	}
	one := col.FindOne(ctx, selector)
	result := bson.M{}
	one.Decode(&result)
	log.Println("===============")
	log.Println(result)
}

func (self *PlaygroundSuite) TestUpdateOne() {
	ctx := context.Background()
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
	result, err := col.UpdateOne(ctx, selector, updater)
	if err != nil {
		panic(err)
	}
	log.Println("===============")
	log.Println(result)
}
