package main

import (
	"os"

	"github.com/spf13/viper"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/server/jp"
	"cybex-gateway/server/user"
	"cybex-gateway/utils/log"
	"cybex-gateway/worker/cyborder"
	jpworker "cybex-gateway/worker/jp"
	"cybex-gateway/worker/order"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	// init configs
	config.LoadConfig(env)
	logDir := viper.GetString("log.log_dir")
	logLevel := viper.GetString("log.log_level")
	log.InitLog(logDir, logLevel, "")
	model.INITFromViper()

	cyborder.InitNode()
	cyborder.InitAsset()
	// start worker and server
	go cyborder.HandleWorker(5)
	go cyborder.BlockRead()

	go order.HandleWorker(5)
	go user.StartServer()
	go jpworker.HandleWorker(5)
	jp.StartServer()
}
