package main

import (
	"blogrpc/openapi/business/controller"
	"blogrpc/openapi/business/server"
	"fmt"
	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
	conf "github.com/spf13/viper"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	env  = flag.String("env", "local", "the running env, development or production")
	port = flag.String("port", "9091", "the listen port, default 9091")
	host = flag.String("host", "0.0.0.0", "the http server listen ip")
)

func main() {
	//// to avoid sighup signal turn down the whole server
	//signal.Ignore(syscall.SIGHUP)
	//
	//loadConfig()
	//
	//conf.Set("logger-level", "debug")
	//os.Setenv("MONGO_MASTER_DSN", "mongodb://root:root@mongo:27017/portal-master?authSource=admin")
	//os.Setenv("MONGO_MASTER_REPLSET", "none")
	//os.Setenv("CACHE_HOST", "redis")
	//os.Setenv("CACHE_PORT", "6379")
	//os.Setenv("RESQUE_HOST", "redis")
	//os.Setenv("RESQUE_PORT", "6379")
	//conf.Set("extension-redis", map[string]interface{}{
	//	"db":        "1", // 注意 redis 使用的是哪个 db，每个服务需要一样才能取到对应的值
	//	"resque-db": "2",
	//})
	//coreLog.InitLogger(conf.GetString("logger-level"), *env, "openapi")
	//
	//log.Printf("API server starts at env: %s, version: %s", *env, util.GetVersion())
	//s := server.NewApiServer(*host, *port, util.GetVersion(), *env)
	//err := s.Boot()
	//if nil != err {
	//	log.Fatalf("Failed to boot server with env: %s, error: %v", *env, err)
	//}
	//s.Run()

	mux := server.NewServeMux()

	engin := gin.New()
	engin.Use(gin.Recovery())
	engin.GET("/ping", func(c *gin.Context) {
		hostname, _ := os.Hostname()
		c.JSON(http.StatusOK, map[string]interface{}{
			"hostname": hostname,
			"ip":       getIp(),
		})
	})
	engin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	engin.GET("/accessToken", controller.AccessTokenHandler)

	engin.Use(controller.Auth)

	engin.Any("/v1/*rest", func(c *gin.Context) {
		mux.ServeHTTP(c.Writer, c.Request)
	})

	engin.Run(":9091")
	graceful.Run(":8080", 4*time.Second, engin)
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

	// 调试时设置绝对路径，不然找不到文件
	conf.SetConfigFile(fmt.Sprintf(confFormat, "/home/user/GolandProjects/blogrpc/openapi/business/conf", *env))
	//conf.SetConfigFile(fmt.Sprintf(confFormat, "./conf", *env))
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

func getIp() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
