package service

import (
	"blogrpc/proto/hello"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func (HelloService) Hello(ctx context.Context, req *hello.StringMessage) (*hello.StringMessage, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	fmt.Println(ok)
	fmt.Println(md)

	return &hello.StringMessage{Value: "hello"}, nil
}
