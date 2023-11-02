package main

import (
	"blogrpc/core/util"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

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

	fmt.Println(util.GetIp())
}
