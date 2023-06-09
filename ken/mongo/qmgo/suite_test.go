package qmgo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PlaygroundSuite2 struct {
	BaseFeatureSuite2
}

func TestPlaygroundSuite2(t *testing.T) {
	mallSuit := new(PlaygroundSuite2)
	mallSuit.Suite = new(suite.Suite)
	suite.Run(t, mallSuit)
}

//func (self *PlaygroundSuite2) log(args ...interface{}) {
//	log.Println(args...)
//}
