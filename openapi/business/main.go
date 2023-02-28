package main

import (
	"blogrpc/openapi/business/controller"
	"blogrpc/openapi/business/server"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

var (
	env = flag.String("env", "local", "the running env, development or production")
)

func main() {
	signal.Ignored(syscall.SIGHUP)

	log.Printf("API server starts at env: %s, version: %s", *env, "v1")

	mux := server.NewServeMux()

	engin := gin.New()
	engin.Use(gin.Recovery())
	engin.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	engin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	engin.GET("/accessToken", controller.GetAccessTokenHandler)

	engin.Use(controller.Auth)

	engin.Any("/v1/*rest", func(c *gin.Context) {
		mux.ServeHTTP(c.Writer, c.Request)
	})

	engin.Run(":9091")
	//graceful.Run(":8080", 4*time.Second, engin)
}
