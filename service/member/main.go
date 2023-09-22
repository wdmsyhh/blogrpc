package main

import (
	"blogrpc/core/constant"
	"blogrpc/core/extension"
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

	// 本地调试的时候使用
	setEnv()

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
	os.Setenv("MONGO_MASTER_REPLSET", "none")
}
