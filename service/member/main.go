package main

import (
	"blogrpc/proto/member"
	"blogrpc/service/member/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.Lshortfile)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	member.RegisterMemberServiceServer(server, &service.MemberService{})
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
