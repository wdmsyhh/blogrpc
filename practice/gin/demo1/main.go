package main

import "github.com/gin-gonic/gin"

func main() {
	e := gin.Default()
	e.POST("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	e.Run("127.0.0.1:8080")
}
