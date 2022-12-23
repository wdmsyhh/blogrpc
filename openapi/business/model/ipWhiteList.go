package model

import (
	"fmt"

	"blogrpc/core/extension"
	"blogrpc/openapi/business/util"

	"blogrpc/core/extension/bson"
	gin "github.com/gin-gonic/gin"
)

const (
	C_IPWHITELIST = "ipWhiteList"
)

type IPWhiteList struct {
	Id        bson.ObjectId `bson:"_id"`
	IP        string        `bson:"ip"`
	AccountId bson.ObjectId `bson:"accountId"`
	BaseModel BaseModel     `bson:"-"`
}

func (self *IPWhiteList) GetByIpAndAccountId(c *gin.Context, accountId bson.ObjectId, ip string) {
	condition := bson.M{
		"ip":        ip,
		"accountId": accountId,
	}
	err := extension.DBRepository.FindOne(c, C_IPWHITELIST, condition, self)
	err = self.BaseModel.filterError(err)
	if err != nil {
		msg := fmt.Sprintf("Error find blacklist by condition: %v, error: %v", condition, err)
		panic(util.NewApiError(util.ModelDBError, msg, err))
	}
}

func IsIpInWhiteList(c *gin.Context, accountId, ip string) bool {
	ipWhiteList := new(IPWhiteList)
	ipWhiteList.GetByIpAndAccountId(c, bson.ObjectIdHex(accountId), ip)
	if !ipWhiteList.Id.Valid() {
		return false
	}

	return true
}
