package server

import (
	"blogrpc/proto"
	"context"
	"log"
	"net/http"
)

func NewServeMux() http.Handler {
	ctx := context.Background()
	log.Println("Register blogrpc service")
	serveMux, err := proto.NewGateway(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Register blogrpc service end")
	return serveMux
}
