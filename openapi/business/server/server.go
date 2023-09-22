package server

import (
	"blogrpc/openapi/business/middleware"
	"blogrpc/openapi/business/router"
	"blogrpc/proto"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	default404Body = "Resource not found"
	default405Body = "Method not allowed"
)

type ApiServer struct {
	Ip           string
	Port         string
	Version      string
	Env          string
	Engine       *gin.Engine
	ServerStatus *middleware.Status
	Mux          http.Handler
}

func NewApiServer(ip, port, version, env string) *ApiServer {
	if env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.HandleMethodNotAllowed = true

	return &ApiServer{
		Engine:  engine,
		Ip:      ip,
		Port:    port,
		Version: version,
		Env:     env,
		Mux:     newServeMux(),
	}
}

func NewServeMux() http.Handler {
	return newServeMux()
}

func newServeMux() http.Handler {
	ctx := context.Background()
	log.Println("Register blogrpc service")
	serveMux, err := proto.NewGateway(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Register blogrpc service end")
	return serveMux
}

func (self *ApiServer) Boot() error {
	return Bootstrap(self)
}

func (self *ApiServer) DebugMode() bool {
	return self.Env == "local"
}

func (self *ApiServer) DevMode() bool {
	return self.Env == "dev"
}

func (self *ApiServer) isTestMode() bool {
	return !strings.Contains(self.Env, "production")
}

func (self *ApiServer) Run() {
	defer Unload()
	addr := fmt.Sprintf("%s:%s", self.Ip, self.Port)
	graceful.Run(addr, 4*time.Second, self.Engine)
}

func (self *ApiServer) GenerateV2HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		self.Mux.ServeHTTP(c.Writer, c.Request)
	}
}

func (self *ApiServer) GetServerStatus(c *gin.Context) {
	c.JSON(http.StatusOK, self.ServerStatus)
}

func (self *ApiServer) Ping(c *gin.Context) {
	if !self.ServerStatus.DBReachable {
		c.JSON(http.StatusServiceUnavailable, self.ServerStatus)
	} else {
		c.Status(http.StatusOK)
	}
}

func (self *ApiServer) NotFound(c *gin.Context) {
	httpMethod := strings.ToUpper(c.Request.Method)
	path := c.Request.URL.Path

	pass, params, node := router.Find(httpMethod, path)
	if pass && node != nil && node.Action != nil {
		for k, v := range params {
			c.Params = append(c.Params, gin.Param{
				Key:   k,
				Value: v,
			})
		}
		handler := node.Action.GenerateHandler()
		handler(c)
		return
	}

	if self.Engine.HandleMethodNotAllowed {
		for method, tree := range router.Trees {
			if method != httpMethod {
				pass, _, _ := tree.Find(path)
				if pass {
					self.MethodNotAllowed(c)
					return
				}
			}
		}
	}

	c.JSON(http.StatusNotFound, map[string]string{"message": default404Body})
	return
}

func (self *ApiServer) MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, map[string]string{"message": default405Body})
}
