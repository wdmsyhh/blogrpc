package mgo

import (
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/suite"
)

var (
	err        error
	mgoSession *mgo.Session
)

type BaseFeatureSuite struct {
	*suite.Suite
}

func (suite *BaseFeatureSuite) SetupTest() {
	mgoSession, err = mgo.Dial("mongodb://root:root@127.0.0.1:27012/?authSource=admin")
	if err != nil {
		panic(err)
	}
}
