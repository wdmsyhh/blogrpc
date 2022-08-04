package server

import (
	"blogrpc/proto"
	"context"
	"log"
	"net/http"
)

func NewServeMux() http.Handler {
	ctx := context.Background()
	log.Println("Register blogrpc services")
	serveMux, err := proto.NewGateway(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Register blogrpc services end")
	return serveMux
}
