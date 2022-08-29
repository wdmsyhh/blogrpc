package share

import (
	"blogrpc/proto/hello"
	"google.golang.org/grpc"
	"log"
)

func GetHelloClient() hello.HelloServiceClient {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cann't connect grpc hello: %v", err)
	}

	client := hello.NewHelloServiceClient(conn)
	return client
}
