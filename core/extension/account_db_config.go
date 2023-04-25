package extension

import (
	"blogrpc/core/extension/bson"
	"context"
)

const (
	C_ACCOUNT_DB_CONFIG = "accountDBConfig"
)

var (
	CAccountDBConfig = &AccountDBConfig{}
)

type AccountDBConfig struct {
	Id        bson.ObjectId          `bson:"_id"`
	Title     string                 `bson:"title"`
	DSN       string                 `bson:"dsn"`
	Options   map[string]interface{} `bson:"options"`
	AccountId bson.ObjectId          `bson:"accountId"`
}

func (*AccountDBConfig) Get(ctx context.Context, accountId string) *AccountDBConfig {
	db := new(AccountDBConfig)
	selector := bson.M{
		"accountId": bson.ObjectIdHex(accountId),
	}
	DBRepository.FindOne(ctx, C_ACCOUNT_DB_CONFIG, selector, db)

	return db
}

func (*AccountDBConfig) GetAll(ctx context.Context) []AccountDBConfig {
	var dbs []AccountDBConfig
	DBRepository.FindAll(ctx, C_ACCOUNT_DB_CONFIG, nil, nil, 0, &dbs)

	return dbs
}

func (*AccountDBConfig) GetByAccountIds(ctx context.Context, accountIds []bson.ObjectId) []AccountDBConfig {
	var dbs []AccountDBConfig
	selector := bson.M{
		"accountId": bson.M{"$in": accountIds},
	}

	DBRepository.FindAll(ctx, C_ACCOUNT_DB_CONFIG, selector, []string{}, 0, &dbs)
	return dbs
}

func (m *AccountDBConfig) Host() string {
	return getMgoConnectString(m.DSN, m.Options)
}
