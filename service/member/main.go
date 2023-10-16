package main

import (
	"blogrpc/core/constant"
	"blogrpc/core/extension"
	"blogrpc/core/util"
	"blogrpc/proto/member"
	"blogrpc/service/member/service"
	flag "github.com/spf13/pflag"
	conf "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	lis, err := net.Listen("tcp", ":"+constant.SERVICE_MEMBER_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if util.IsRunningInContainer() {
		log.Println("====IsRunningInContainer===")
		log.Println(os.Getenv("MONGO_MASTER_DSN"))
		log.Println(os.Getenv("MONGO_MASTER_REPLSET"))
	} else {
		// 本地调试的时候使用
		setEnv()
	}

	extension.LoadExtensionsByName([]string{"mgo"}, *env == Local)

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

	os.Setenv("MONGO_MASTER_DSN", "mongodb://root:root@localhost:27012/portal-master?authSource=admin")
	os.Setenv("MONGO_MASTER_REPLSET", "rs0")

	os.Setenv("MYSQL_MASTER_DSN", "root:root123@tcp(localhost:3306)/portal_master?charset=utf8mb4&parseTime=True&loc=Local")
}
