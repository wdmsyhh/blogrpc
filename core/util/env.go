package util

import (
	"fmt"
	"os"
)

func GetMongoAppName() string {
	return fmt.Sprintf("%s.%s", os.Getenv("K8S_SERVICE_NAME"), os.Getenv("K8S_SERVICE_NAMESPACE"))
}

func GetMongoMasterDsn() string {
	return os.Getenv("MONGO_MASTER_DSN")
}

func GetMongoMasterReplset() string {
	return os.Getenv("MONGO_MASTER_REPLSET")
}

func GetMysqlMasterDsn() string {
	return os.Getenv("MYSQL_MASTER_DSN")
}

func GetCacheHost() string {
	return os.Getenv("CACHE_HOST")
}

func GetCachePort() string {
	return os.Getenv("CACHE_PORT")
}

func GetCachePassword() string {
	return os.Getenv("CACHE_PASSWORD")
}

func GetResqueHost() string {
	return os.Getenv("RESQUE_HOST")
}

func GetResquePort() string {
	return os.Getenv("RESQUE_PORT")
}

func GetResquePassword() string {
	return os.Getenv("RESQUE_PASSWORD")
}
