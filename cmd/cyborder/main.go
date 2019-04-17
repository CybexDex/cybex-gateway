package main

import (
	"os"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"bitbucket.org/woyoutlz/bbb-gateway/worker/cyborder"
	"github.com/spf13/viper"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)

	logDir := viper.GetString("log.log_dir")
	logLevel := viper.GetString("log.log_level")
	log.InitLog(logDir, logLevel, "[bbb]")

	model.INITFromViper()
	cyborder.InitNode()
	cyborder.InitAsset()
	go cyborder.HandleWorker(5)
	go cyborder.BlockRead()
	select {}
	// fmt.Println(s)
}
