package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"testing"
	"time"

	registry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestEtcdTls(t *testing.T) {
	tlsCert := ""
	tlsKey := ""
	tlsCa := ""

	ct, err := tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey))
	if err != nil {
		log.Panicf("[etcd] parses etcd tls cert error:%s", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(tlsCa))

	cfg := etcd.Config{
		Endpoints:            []string{"10.64.144.105:2379"},
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    5 * time.Second, // 每5秒ping一次服务器
		DialKeepAliveTimeout: time.Second,     // 1秒没有返回则代表故障
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
	}

	cfg.TLS = &tls.Config{
		Certificates: []tls.Certificate{ct},
		RootCAs:      pool,
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

	reg := registry.New(cli, registry.Namespace(fmt.Sprintf("/%s/%s", "namespace1", "dev")))

	fmt.Println(reg == nil)

}
