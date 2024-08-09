package main

import (
	"context"
	"fmt"
	"log"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func main() {
	cfg := etcd.Config{
		Endpoints:            []string{"127.0.0.1:2379"},
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    5 * time.Second, // 每5秒ping一次服务器
		DialKeepAliveTimeout: time.Second,     // 1秒没有返回则代表故障
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
	}

	cli, err := etcd.New(cfg)
	if err != nil {
		log.Panicf("[etcd] Failed to etcd:%s", err)
	}

	list, err := cli.MemberList(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
	userList, err := cli.UserList(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(userList)
	result, err := cli.Get(context.Background(), "", etcd.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range result.Kvs {
		fmt.Println(kv)
	}
}
