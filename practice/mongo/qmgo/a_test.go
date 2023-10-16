package qmgo

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
)

func (self *PlaygroundSuite2) TestInsertOne() {
	ctx := context.Background()
	doc := primitive.M{
		"_id":  primitive.NewObjectID(),
		"name": "hh",
	}
	res, err := qmgoClient.InsertOne(ctx, &doc)
	if err != nil {
		fmt.Println(err.Error())
	}
	typ := reflect.TypeOf(res.InsertedID)
	fmt.Println(typ)
	log.Println("===============")
	log.Println(res)
}

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
