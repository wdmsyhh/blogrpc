package client

import (
	pb_member "blogrpc/proto/member"
	"strings"
)

func getMemberServiceClient(cc *ClientConn) interface{} {
	client := pb_member.NewMemberServiceClient(cc.Conn)
	return &client
}

func GetClientByFuncName(funcName string) func(cc *ClientConn) interface{} {
	arr := strings.Split(funcName, ".")
	if len(arr) < 2 {
		return nil
	}

	switch arr[0] {
	case "MemberService":
		return getMemberServiceClient
	}

	return nil
}
