package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserInfo struct {
	Id     uint
	Name   string
	Gender string
	Hobby  string
}

func main() {
	// 连接数据库
	dsn := "root:root123@tcp(127.0.0.1:3306)/portal-master?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//自动迁移
	db.AutoMigrate(&UserInfo{})
	u1 := UserInfo{Id: 1, Name: "张三", Gender: "男", Hobby: "学习"}
	db.Create(&u1) //创建
}
