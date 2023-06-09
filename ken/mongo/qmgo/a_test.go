package qmgo

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (self *PlaygroundSuite2) TestGetOne() {
	ctx := context.Background()
	id, _ := primitive.ObjectIDFromHex("6480319cfbcabc003bd46004")
	selector := primitive.M{
		"_id": id,
	}
	find := qmgoClient.Find(ctx, selector)
	result := bson.M{}
	err := find.One(&result)
	if err != nil {
		panic(err)
	}
	log.Println("===============")
	log.Println(result)
}

func (self *PlaygroundSuite2) TestGetOnev2() {
	ctx := context.Background()
	id := bson.ObjectIdHex("6480319cfbcabc003bd46004")
	selector := bson.M{
		"_id": id,
	}
	find := qmgoClient.Find(ctx, selector)
	result := bson.M{}
	err := find.One(&result)
	if err != nil {
		panic(err)
	}
	log.Println("===============")
	log.Println(result)
}

func (self *PlaygroundSuite2) TestUpdateOne() {
	ctx := context.Background()
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
	err := qmgoClient.UpdateOne(ctx, selector, updater)
	if err != nil {
		panic(err)
	}
}
