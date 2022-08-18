package service

import "blogrpc/proto/hello"

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}
