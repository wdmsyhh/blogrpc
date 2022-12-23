package model

import "blogrpc/core/extension/bson"

const (
	QueryTimeFormat = "2006-01-02"
)

type BaseModel struct {
}

func (self *BaseModel) filterError(err error) error {
	if err == bson.ErrNotFound {
		return nil
	}
	return err
}
