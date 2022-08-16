package service

import (
	"blogrpc/core/extension"
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/services/share"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (MemberService) GetMember(ctx context.Context, req *member.GetMemberRequest) (*member.GetMemberResponse, error) {

	resp := &member.GetMemberResponse{}

	if req.Id == "aaa" {
		resp = &member.GetMemberResponse{
			Id:   "aaa",
			Name: "小明",
			Age:  20,
		}
	}

	helloResp, err := share.GetHelloClient().SayHello(ctx, &hello.StringMessage{Value: "Hello a"})
	if err != nil {
		return nil, err
	}

	user := struct {
		Id     primitive.ObjectID `bson:"_id"`
		Name   string             `bson:"name"`
		Age    uint16             `bson:"age"`
		Weight uint32             `bson:"weight"`
	}{}
	id, _ := primitive.ObjectIDFromHex("62fb053a54099e3978aeb655")
	selector := bson.M{
		"_id": id,
	}
	err = extension.DBRepository.FindOne(ctx, "", selector, &user)
	if err != nil {
		return nil, err
	}

	resp.Name = helloResp.Value + user.Name
	resp.Age = int64(user.Age)
	resp.Id = user.Id.Hex()

	return resp, nil
}
