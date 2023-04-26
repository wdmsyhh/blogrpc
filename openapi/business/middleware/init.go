package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	Name() string
	InitWithConf(conf map[string]interface{}, debug bool) error
	MiddlewareFunc() gin.HandlerFunc
}

var (
	//name and the middleware instance
	middlewares map[string]IMiddleware = make(map[string]IMiddleware)
)

func GetMiddlewares() map[string]IMiddleware {
	return middlewares
}

func enableMiddleware(m IMiddleware) {
	if nil != m {
		name := m.Name()
		middlewares[name] = m
	} else {
		log.Panicln("Failed to enable middleware, middleware is nil.")
	}
}
