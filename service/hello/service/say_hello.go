package service

import (
	"blogrpc/proto/hello"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"os"
)

func (HelloService) SayHello(ctx context.Context, req *hello.StringMessage) (*hello.StringMessage, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	fmt.Println(ok)
	fmt.Println(md)

	hostname, _ := os.Hostname()
	resp := &hello.StringMessage{
		Value:    req.Value,
		Hostname: hostname,
	}

	return resp, nil
}
