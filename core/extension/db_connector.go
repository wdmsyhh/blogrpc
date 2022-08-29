package extension

import (
	"context"
	"github.com/qiniu/qmgo"
)

var (
	TenantDBConnector *tenantDBConnector = nil
)

func init() {
	TenantDBConnector = &tenantDBConnector{}
}

type tenantDBConnector struct {
}

type DBConnector interface {
	GetClient(ctx context.Context, collection string) *qmgo.QmgoClient
}

func (*tenantDBConnector) GetClient(ctx context.Context, collection string) *qmgo.QmgoClient {
	// 如果是在容器内链接另一个 mongo 容器，需要使用 mongo 容器的内部端口
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27012", Database: "class", Coll: collection, Auth: &qmgo.Credential{
		Username: "admin",
		Password: "123456",
	},
	})
	if err != nil {
		panic(err)
	}
	return cli
}
