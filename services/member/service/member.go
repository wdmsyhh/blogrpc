package service

import (
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/services/share"
	"context"
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

	helloResp, err := share.GetHelloClient().Hello(ctx, &hello.StringMessage{})
	if err != nil {
		return nil, err
	}

	resp.Name = helloResp.Value

	return resp, nil
}
