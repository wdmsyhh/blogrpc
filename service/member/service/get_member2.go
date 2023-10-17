package service

import (
	"blogrpc/core/client"
	"blogrpc/core/extension/mysql"
	"blogrpc/core/util"
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"blogrpc/service/member/model"
	"context"
	"fmt"
	"os"
)

func (MemberService) GetMember2(ctx context.Context, req *member.GetMemberRequest2) (*member.GetMemberResponse2, error) {

	resp := &member.GetMemberResponse2{}

	hostname, _ := os.Hostname()
	if req.Id == 1111 {
		resp = &member.GetMemberResponse2{
			Id:      req.Id,
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

	if req.Id == 2222 {
		resp = &member.GetMemberResponse2{
			Id:      req.Id,
			Name:    "小明",
			Age:     20,
			Service: "member-" + hostname + "-" + util.GetIp() + ";" + helloResp.Service,
		}
		return resp, nil
	}

	var dbMember model.TMember
	tx := mysql.DB.First(&dbMember, req.Id)

	fmt.Println(tx)

	resp.Name = helloResp.Value + dbMember.Name
	resp.Age = dbMember.Age
	resp.Id = int64(dbMember.ID)
	resp.Service = "member-" + hostname + "-" + util.GetIp() + ";" + helloResp.Service
	//resp.Name = helloResp.Value

	return resp, nil
}
