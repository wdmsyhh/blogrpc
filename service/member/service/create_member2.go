package service

import (
	"blogrpc/core/extension/mysql"
	pb_member "blogrpc/proto/member"
	"blogrpc/service/member/model"
	"context"
	"fmt"
)

func (MemberService) CreateMember2(ctx context.Context, req *pb_member.CreateMemberRequest2) (*pb_member.CreateMemberResponse2, error) {

	member := &model.TMember{
		Name: req.Name,
		Age:  req.Age,
	}

	mysql.DB.AutoMigrate(&model.TMember{})

	tx := mysql.DB.Create(member)

	fmt.Println("=====tx====", tx)

	return &pb_member.CreateMemberResponse2{
		Id:   int64(member.ID),
		Name: member.Name,
		Age:  member.Age,
	}, nil
}
