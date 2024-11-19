package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := gin.Default()
	e.POST("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	user := &User{
		DB:          "mysql",
		RedisClient: NewTestRedis(),
	}
	e.POST("/user", user.GetUser)
	e.Run("127.0.0.1:8080")
}

type User struct {
	DB          string `json:"db"`
	RedisClient *redis.Client
}

type GetUserRequest struct {
	Id int64 `json:"id"`
}

func (u *User) GetUser(ctx *gin.Context) {
	req := &GetUserRequest{}
	_ = ctx.BindJSON(req)

	name := ctx.Query("name")
	u.RedisClient.Set(context.Background(), "name", name, 0)

	value := u.RedisClient.Get(context.Background(), "name").Val()
	fmt.Println(value)

	if req.Id == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"id":   1,
			"name": "小明",
			"age":  18,
			"db":   u.DB,
		})
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{"code": "UserNotExist"})
	return
}

func NewTestRedis() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}
