package main

import (
	"os"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
	"cybex-gateway/worker/cyborder"

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
	// go cyborder.HandleWorker(5)
	// go cyborder.BlockRead()
	cyborder.HandleBlockNum(9024999)
	// select {}
	// fmt.Println(s)
}
