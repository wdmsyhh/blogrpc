package client

import (
	"log"

	"blogrpc/core/constant"
	"blogrpc/proto/hello"
	"google.golang.org/grpc"
)

func GetHelloServiceClient() hello.HelloServiceClient {
	conn, err := grpc.Dial(constant.ServiceHelloHost+":"+constant.ServiceHelloPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cann't connect grpc hello: %v", err)
	}
	client := hello.NewHelloServiceClient(conn)
	return client
}
