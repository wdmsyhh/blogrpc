package mongo

import (
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type PlaygroundSuite struct {
	BaseFeatureSuite
}

func TestPlaygroundSuite(t *testing.T) {
	mallSuit := new(PlaygroundSuite)
	mallSuit.Suite = new(suite.Suite)
	suite.Run(t, mallSuit)
}

func (self *PlaygroundSuite) log(args ...interface{}) {
	log.Println(args...)
}
