package main

import (
	coreLog "blogrpc/core/log"
	"blogrpc/openapi/business/server"
	"blogrpc/openapi/business/util"
	"fmt"
	flag "github.com/spf13/pflag"
	conf "github.com/spf13/viper"
	"log"
	"os/signal"
	"syscall"
)

var (
	env  = flag.String("env", "local", "the running env, development or production")
	port = flag.String("port", "9091", "the listen port, default 9091")
	host = flag.String("host", "0.0.0.0", "the http server listen ip")
)

func main() {
	// to avoid sighup signal turn down the whole server
	signal.Ignore(syscall.SIGHUP)

	loadConfig()
	coreLog.InitLogger(
		conf.GetString("logger-level"),
		*env,
		"openapi",
	)

	log.Printf("API server starts at env: %s, version: %s", *env, util.GetVersion())
	s := server.NewApiServer(*host, *port, util.GetVersion(), *env)
	err := s.Boot()
	if nil != err {
		log.Fatalf("Failed to boot server with env: %s, error: %v", *env, err)
	}
	s.Run()

	//mux := server.NewServeMux()
	//
	//engin := gin.New()
	//engin.Use(gin.Recovery())
	//engin.GET("/ping", func(c *gin.Context) {
	//	c.Status(http.StatusOK)
	//})
	//engin.GET("/health", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    http.StatusOK,
	//		"success": true,
	//	})
	//})
	//
	//engin.GET("/accessToken", controller.GetAccessTokenHandler)
	//
	//engin.Use(controller.Auth)
	//
	//engin.Any("/v1/*rest", func(c *gin.Context) {
	//	mux.ServeHTTP(c.Writer, c.Request)
	//})
	//
	//engin.Run(":9091")
	//graceful.Run(":8080", 4*time.Second, engin)
}

func loadConfig() {
	// Parse the ENV, for example OMNIAPI_MODE
	conf.SetEnvPrefix("omniapi")
	conf.AutomaticEnv()

	conf.BindPFlag("env", flag.Lookup("env"))
	conf.BindPFlag("port", flag.Lookup("port"))
	conf.BindPFlag("host", flag.Lookup("host"))
	// Parse configurations
	flag.Parse()

	confFormat := "%s/%s.toml"

	// read common config
	conf.SetConfigFile(fmt.Sprintf(confFormat, "./conf", *env))
	err := conf.MergeInConfig()
	if err != nil {
		log.Println(err)
	}

	// Add global config to conf
	conf.Set("env", *env)
	conf.Set("addr", fmt.Sprintf("%s:%s", *host, *port))
	conf.Set("service", "openapi")
	log.Printf("Configuration loaded from: %s", conf.ConfigFileUsed())
}
