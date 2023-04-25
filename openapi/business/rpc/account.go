package mairpc

import (
	"blogrpc/core/client"
	"blogrpc/openapi/business/util"
	pb_account "blogrpc/proto/account"
	"blogrpc/proto/common/request"
	"github.com/gin-gonic/gin"
)

var CAccount = &Account{}

type Account struct{}

func (*Account) GetSessionValidPeriod(c *gin.Context) (*pb_account.SessionValidPeriodResponse, error) {
	resp, err := client.Run(
		"AccountService.GetSessionValidPeriod",
		util.GetGrpcContext(c),
		&request.EmptyRequest{},
	)
	if err != nil {
		return nil, err
	}

	return resp.(*pb_account.SessionValidPeriodResponse), nil
}
