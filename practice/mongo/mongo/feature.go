package mongo

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err    error
	client *mongo.Client
)

type BaseFeatureSuite struct {
	*suite.Suite
}

func (suite *BaseFeatureSuite) SetupTest() {
	ctx := context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27012/?authSource=admin").SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	}))
	if err != nil {
		panic(err)
	}
}
