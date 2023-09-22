package client

import (
	"blogrpc/core/constant"
	"blogrpc/proto/hello"
	"google.golang.org/grpc"
	"log"
)

func GetHelloServiceClient() hello.HelloServiceClient {
	conn, err := grpc.Dial(constant.SERVICE_HELLO_HOST+":"+constant.SERVICE_HELLO_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cann't connect grpc hello: %v", err)
	}
	client := hello.NewHelloServiceClient(conn)
	return client
}
