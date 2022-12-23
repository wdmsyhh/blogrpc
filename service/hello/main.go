package main

import (
	"blogrpc/proto/hello"
	"blogrpc/service/hello/service"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	//go StartGRPCGateway()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	hello.RegisterHelloServiceServer(server, &service.HelloService{})
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func StartGRPCGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}), runtime.WithIncomingHeaderMatcher(func(headerName string) (string, bool) {

		if headerName == "App-Id" {
			return "appid", true
		}

		return "", false
	}))

	err := hello.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, ":8081", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("cann't start grpc gateway: %v", err)
	}

	err = http.ListenAndServe(":8080", mux) // grpc gateway 的
	if err != nil {
		log.Fatalf("cann't listen and serve: %v", err)
	}
}
