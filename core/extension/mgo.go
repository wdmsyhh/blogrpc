package extension

import (
	"context"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DBRepository DatabaseRepository = &mongoDBRepository{
		clientSelector: &proxyClientSelector{},
	}
)

type mongoDBRepository struct {
	clientSelector ClientSelector
}

type ClientSelector interface {
	GetClient(ctx context.Context, collectionName string) *qmgo.QmgoClient
}

type proxyClientSelector struct{}

func (*proxyClientSelector) GetClient(ctx context.Context, collectionName string) *qmgo.QmgoClient {
	return getClient(ctx, collectionName)
}

func getClient(ctx context.Context, collection string) *qmgo.QmgoClient {
	client := TenantDBConnector.GetClient(ctx, collection)
	return client
}

type DatabaseRepository interface {
	FindOne(ctx context.Context, collectionName string, selector bson.M, result interface{}) error
	Insert(ctx context.Context, collectionName string, docs ...interface{}) error
}

func (repo *mongoDBRepository) FindOne(ctx context.Context, collectionName string, selector bson.M, result interface{}) error {
	cli := repo.clientSelector.GetClient(ctx, collectionName)
	defer func() {
		if err := cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	// 传入的 result 不能是双重指针
	return cli.Find(ctx, selector).One(result)
}

func (repo *mongoDBRepository) Insert(ctx context.Context, collectionName string, docs ...interface{}) error {
	cli := repo.clientSelector.GetClient(ctx, collectionName)
	defer func() {
		if err := cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	_, err := cli.InsertMany(ctx, docs)
	return err
}
