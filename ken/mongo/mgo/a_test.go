package mgo

import (
	"github.com/globalsign/mgo/bson"
	"log"
)

func (self *PlaygroundSuite) TestGetOne() {
	col := mgoSession.DB("testdb").C("cola")
	id := bson.ObjectIdHex("6480319cfbcabc003bd46004")
	selector := bson.M{
		"_id": id,
	}
	find := col.Find(selector)
	result := bson.M{}
	err := find.One(&result)
	if err != nil {
		panic(err)
	}
	log.Println("===============")
	log.Println(result)
}

func (self *PlaygroundSuite) TestUpdateOne() {
	col := mgoSession.DB("testdb").C("cola")
	id := bson.ObjectIdHex("6480319cfbcabc003bd46004")
	selector := bson.M{
		"_id":       id,
		"isDeleted": false,
	}
	updater := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}
	err := col.Update(selector, updater)
	if err != nil {
		panic(err)
	}
}
