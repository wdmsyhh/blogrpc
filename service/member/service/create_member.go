package service

import (
	pb_member "blogrpc/proto/member"
	"blogrpc/service/member/model"
	"context"
	"github.com/qiniu/qmgo"
)

func (MemberService) CreateMember(ctx context.Context, req *pb_member.CreateMemberRequest) (*pb_member.CreateMemberResponse, error) {

	member := &model.Member{
		Id:   qmgo.NewObjectID(),
		Name: req.Name,
		Age:  req.Age,
	}

	id, err := member.Create(ctx)
	if err != nil {
		return nil, err
	}

	return &pb_member.CreateMemberResponse{
		Id:   id.Hex(),
		Name: member.Name,
		Age:  member.Age,
	}, nil
}
