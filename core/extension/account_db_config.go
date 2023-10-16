package extension

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	C_ACCOUNT_DB_CONFIG = "accountDBConfig"
)

var (
	CAccountDBConfig = &AccountDBConfig{}
)

type AccountDBConfig struct {
	Id        primitive.ObjectID     `bson:"_id"`
	Title     string                 `bson:"title"`
	DSN       string                 `bson:"dsn"`
	Options   map[string]interface{} `bson:"options"`
	AccountId primitive.ObjectID     `bson:"accountId"`
}

func (*AccountDBConfig) Get(ctx context.Context, accountId string) *AccountDBConfig {
	db := new(AccountDBConfig)
	aid, _ := primitive.ObjectIDFromHex(accountId)
	selector := primitive.M{
		"accountId": aid,
	}
	DBRepository.FindOne(ctx, C_ACCOUNT_DB_CONFIG, selector, db)

	return db
}

func (*AccountDBConfig) GetAll(ctx context.Context) []AccountDBConfig {
	var dbs []AccountDBConfig
	DBRepository.FindAll(ctx, C_ACCOUNT_DB_CONFIG, nil, nil, 0, &dbs)

	return dbs
}

func (*AccountDBConfig) GetByAccountIds(ctx context.Context, accountIds []primitive.ObjectID) []AccountDBConfig {
	var dbs []AccountDBConfig
	selector := primitive.M{
		"accountId": primitive.M{"$in": accountIds},
	}

	DBRepository.FindAll(ctx, C_ACCOUNT_DB_CONFIG, selector, []string{}, 0, &dbs)
	return dbs
}

func (m *AccountDBConfig) Host() string {
	return getMgoConnectString(m.DSN, m.Options)
}
