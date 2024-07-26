package router

import (
	"context"
	"net/http"
	"time"

	"blogrpc/practice/urfavecli/server"
	"github.com/gin-gonic/gin"
)

type Cmd2Router struct {
}

func NewCmd2Router() *Cmd2Router {
	return &Cmd2Router{}
}

func (r *Cmd2Router) Run(ctx context.Context, addr string) error {
	srv := server.New(
		server.Address(addr),
		server.Handler(r.route()),
		server.StopTimeout(3*time.Second),
	)
	return srv.Start(ctx)
}

func (r *Cmd2Router) route() http.Handler {
	engine := gin.Default()

	engine.Any("/heartbeat/check", func(c *gin.Context) {
		c.String(200, "ok")
	})

	engine.Any("/ping", func(context *gin.Context) {
		context.JSON(200, "pong cmd2")
	})

	engine.GET("path1", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"name": "小明cmd2",
			"age":  18,
		})
	})

	return engine
}
