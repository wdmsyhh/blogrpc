package model

import (
	"blogrpc/core/extension"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	C_MEMBER = "member"
)

var (
	CMember = &Member{}
)

type Member struct {
	Id        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	IsDeleted bool               `bson:"isDeleted"`
	Name      string             `bson:"name"`
	Age       int64              `bson:"age"`
}

func (*Member) GetById(ctx context.Context, id primitive.ObjectID) (Member, error) {
	result := Member{}
	condition := bson.M{
		"_id": id,
	}
	err := extension.DBRepository.FindOne(ctx, C_MEMBER, condition, &result)
	return result, err
}

func (m *Member) Create(ctx context.Context) error {
	m.CreatedAt, m.UpdatedAt = time.Now(), time.Now()
	err := extension.DBRepository.Insert(ctx, C_MEMBER, m)
	return err
}
