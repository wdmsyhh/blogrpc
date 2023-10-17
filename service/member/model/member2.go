package model

import "gorm.io/gorm"

type TMember struct {
	gorm.Model
	Name string
	Age  int64
}
