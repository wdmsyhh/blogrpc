package mysql

var (
	MysqlConnector *mysqlConnector = nil
)

type mysqlConnector struct {
	conf map[string]interface{}
}
