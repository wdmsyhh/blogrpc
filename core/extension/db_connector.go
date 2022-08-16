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
	GetClient(ctx context.Context) *qmgo.QmgoClient
}

func (*tenantDBConnector) GetClient(ctx context.Context) *qmgo.QmgoClient {
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27011", Database: "class", Coll: "user", Auth: &qmgo.Credential{
		Username: "admin",
		Password: "123456",
	},
	})
	if err != nil {
		panic(err)
	}
	return cli
}
