package service

import (
	"context"
	"os"

	"blogrpc/core/client"
	"blogrpc/core/errors"
	"blogrpc/core/util"
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/service/member/codes"
	"blogrpc/service/member/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (MemberService) GetMember(ctx context.Context, req *member.GetMemberRequest) (*member.GetMemberResponse, error) {

	if req.Id == "error" {
		return nil, errors.NewInvalidArgumentError("memberId")
	}
	if req.Id == "error1" {
		return nil, errors.NewInvalidArgumentErrorWithMessage("memberId", "错误了")
	}
	if req.Id == "error2" {
		return nil, errors.NewNotExistsError("member")
	}
	if req.Id == "error3" {
		return nil, codes.NewError(codes.MemberNotFound)
	}
	if req.Id == "error4" {
		return nil, status.Errorf(grpc_codes.InvalidArgument, "member not found")
	}
	if req.Id == "error5" {
		return nil, errors.NewInternal("ssss")
	}

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
