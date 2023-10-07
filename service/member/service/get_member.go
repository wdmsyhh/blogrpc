package service

import (
	"blogrpc/core/client"
	"blogrpc/core/util"
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/service/member/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

func (MemberService) GetMember(ctx context.Context, req *member.GetMemberRequest) (*member.GetMemberResponse, error) {

	resp := &member.GetMemberResponse{}

	hostname, _ := os.Hostname()
	if req.Id == "aaa" {
		resp = &member.GetMemberResponse{
			Id:      "aaa",
			Name:    "小明",
			Age:     20,
			Service: "member-" + hostname + "-" + util.GetIp(),
		}
		return resp, nil
	}

	helloResp, err := client.GetHelloServiceClient().SayHello(ctx, &hello.StringMessage{Value: "Hello a"})
	if err != nil {
		return nil, err
	}

	if req.Id == "bbb" {
		resp = &member.GetMemberResponse{
			Id:      "aaa",
			Name:    "小明",
			Age:     20,
			Service: "member-" + hostname + "-" + util.GetIp() + ";" + helloResp.Service,
		}
		return resp, nil
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
	resp.Service = "member-" + hostname + "-" + util.GetIp() + ";" + helloResp.Service
	//resp.Name = helloResp.Value

	return resp, nil
}
