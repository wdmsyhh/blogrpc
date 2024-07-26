package router

import (
	"context"
	"net/http"
	"time"

	"blogrpc/practice/urfavecli/server"
	"github.com/gin-gonic/gin"
)

type Cmd1Router struct {
}

func NewCmd1Router() *Cmd1Router {
	return &Cmd1Router{}
}

func (r *Cmd1Router) Run(ctx context.Context, addr string) error {
	srv := server.New(
		server.Address(addr),
		server.Handler(r.route()),
		server.StopTimeout(3*time.Second),
	)
	return srv.Start(ctx)
}

func (r *Cmd1Router) route() http.Handler {
	engine := gin.Default()

	engine.Any("/heartbeat/check", func(c *gin.Context) {
		c.String(200, "ok")
	})

	engine.Any("/ping", func(context *gin.Context) {
		context.JSON(200, "pong cmd1")
	})

	engine.GET("path1", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"name": "小明cmd1",
			"age":  18,
		})
	})

	return engine
}
