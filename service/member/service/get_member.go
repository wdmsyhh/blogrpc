package service

import (
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/service/member/model"
	"blogrpc/service/share"
	"context"
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
		return resp, nil
	}

	helloResp, err := share.GetHelloClient().SayHello(ctx, &hello.StringMessage{Value: "Hello a"})
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	dbMember, err := model.CMember.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	resp.Name = helloResp.Value + dbMember.Name
	resp.Age = dbMember.Age
	resp.Id = dbMember.Id.Hex()

	return resp, nil
}
