package mysql

import (
	"blogrpc/core/extension"
	"blogrpc/core/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlConnector *mysqlConnector = nil

	DB *gorm.DB
)

func init() {
	MysqlConnector = &mysqlConnector{
		conf:     make(map[string]interface{}),
		dbClient: nil,
	}
	extension.RegisterExtension(MysqlConnector)
}

type mysqlConnector struct {
	conf     map[string]interface{}
	dbClient *gorm.DB
}

func (*mysqlConnector) Name() string {
	return "mysql"
}

func (m *mysqlConnector) InitWithConf(conf map[string]interface{}, debug bool) error {
	db, err := gorm.Open(mysql.Open(util.GetMysqlMasterDsn()), &gorm.Config{})
	if err != nil {
		return err
	}
	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(100) // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

	m.dbClient = db
	DB = db
	return nil
}

func (m *mysqlConnector) Close() {
	// todo
}
