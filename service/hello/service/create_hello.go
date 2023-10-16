package service

import (
	"blogrpc/proto/common/response"
	"blogrpc/proto/hello"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func (HelloService) CreateHello(ctx context.Context, req *hello.StringMessage) (*response.EmptyResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	fmt.Println(ok)
	fmt.Println(md)

	return &response.EmptyResponse{}, nil
}
