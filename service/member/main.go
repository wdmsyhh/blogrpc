package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"blogrpc/core/constant"
	"blogrpc/core/util"
	"blogrpc/proto/member"
	"blogrpc/service/member/service"
	flag "github.com/spf13/pflag"
	conf "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	Local = "local"
)

var (
	env = flag.String("env", Local, "the running environment")
)

func main() {
	log.SetFlags(log.Lshortfile)

	signal.Ignore(syscall.SIGHUP)

	lis, err := net.Listen("tcp", ":"+constant.ServiceMemberPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if util.IsRunningInContainer() {
		setConf()
		log.Println("====IsRunningInContainer===")
		log.Println(os.Getenv("MONGO_MASTER_DSN"))
		log.Println(os.Getenv("MONGO_MASTER_REPLSET"))
		log.Println(os.Getenv("MYSQL_MASTER_DSN"))
		log.Println("confAllKeys:", conf.AllKeys())
		log.Println("logger-level:", conf.Get("logger-level"))
		log.Println("strategy:", conf.Get("strategy"))
	} else {
		// 本地调试的时候使用
		setEnv()
		log.Println("confAllKeys:", conf.AllKeys())
		log.Println("logger-level:", conf.Get("logger-level"))
		log.Println("strategy:", conf.Get("strategy"))
	}

	//extension.LoadExtensionsByName([]string{"mgo", "mysql", "redis"}, *env == Local)

	server := grpc.NewServer()
	member.RegisterMemberServiceServer(server, &service.MemberService{})
	reflection.Register(server)
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func setEnv() {
	conf.Set("logger-level", "debug")
	conf.Set("strategy", "direct")
	conf.Set("extension-redis.db", "1")
	conf.Set("extension-redis.resque-db", "2")
	conf.Set("extension-redis.response-cache-db", "5")

	// mongodb
	os.Setenv("MONGO_MASTER_DSN", "mongodb://root:root@localhost:27012/portal-master?authSource=admin")
	os.Setenv("MONGO_MASTER_REPLSET", "rs0")

	// mysql
	os.Setenv("MYSQL_MASTER_DSN", "root:root123@tcp(localhost:3306)/portal_master?charset=utf8mb4&parseTime=True&loc=Local")

	// redis
	os.Setenv("CACHE_HOST", "localhost")
	os.Setenv("CACHE_PORT", "6379")
	os.Setenv("CACHE_PASSWORD", "root123")
	os.Setenv("RESQUE_HOST", "localhost")
	os.Setenv("RESQUE_PORT", "6379")
	os.Setenv("RESQUE_PASSWORD", "root123")
}

func setConf() {
	confFormat := "%s/%s.toml"

	// read common config
	conf.SetConfigFile(fmt.Sprintf(confFormat, "./conf", *env))
	conf.MergeInConfig()

	// read mairpc common config
	conf.SetConfigFile(fmt.Sprintf(confFormat, "./conf", "common"))
	conf.MergeInConfig()

	conf.Set("service", "MemberService")
}
