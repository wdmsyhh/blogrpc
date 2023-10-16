package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// 注意：设置环境变量时应该大写，不然获取不到，比如 export PRE_ENV1=abc
// 需要在同一个终端执行才行
// 进入到 main 目录执行 go run main.go

func main() {
	viper.SetEnvPrefix("pre")
	viper.AutomaticEnv()

	get := viper.GetString("env1")
	fmt.Println("get:", get, os.Getenv("PRE_ENV1"))

	allSettings := viper.AllSettings()
	fmt.Println(allSettings)
	keys := viper.AllKeys() // 没有列出来环境变量中的 key
	fmt.Println(keys)

	viper.SetConfigFile(fmt.Sprintf("%s/%s.toml", "./conf", "local"))
	viper.MergeInConfig()

	allSettings = viper.AllSettings()
	fmt.Println(allSettings)
	keys = viper.AllKeys()
	fmt.Println(keys)

	fmt.Println(viper.GetString("vod-provider"))
	fmt.Println(viper.GetString("env1"))
}
