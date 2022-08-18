package service

import (
	"blogrpc/proto/common/response"
	"blogrpc/proto/member"
	"blogrpc/service/member/model"
	"context"
	"github.com/qiniu/qmgo"
)

func (MemberService) CreateMember(ctx context.Context, req *member.CreateMemberRequest) (*response.EmptyResponse, error) {

	member := &model.Member{
		Id:   qmgo.NewObjectID(),
		Name: req.Name,
		Age:  req.Age,
	}

	err := member.Create(ctx)
	if err != nil {
		return nil, err
	}

	return &response.EmptyResponse{}, nil
}
