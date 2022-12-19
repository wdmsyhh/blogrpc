package main

import (
	"blogrpc/proto/member"
	"blogrpc/service/member/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile)

	signal.Ignore(syscall.SIGHUP)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	member.RegisterMemberServiceServer(server, &service.MemberService{})
	reflection.Register(server)
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
