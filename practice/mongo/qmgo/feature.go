package qmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/suite"
)

var (
	err        error
	qmgoClient *qmgo.QmgoClient
)

type BaseFeatureSuite2 struct {
	*suite.Suite
}

func (suite *BaseFeatureSuite2) SetupTest() {
	ctx := context.Background()
	qmgoClient, err = qmgo.Open(ctx, &qmgo.Config{
		Uri:      "mongodb://root:root@127.0.0.1:27012",
		Database: "testdb",
		Coll:     "cola",
	})
	if err != nil {
		panic(err)
	}
}
