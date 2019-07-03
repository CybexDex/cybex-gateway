package main

import (
	"os"

	"github.com/spf13/viper"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/server/jpselect"
	"cybex-gateway/server/user"
	"cybex-gateway/utils/log"
	"cybex-gateway/worker/cyborder"
	jpworker "cybex-gateway/worker/jpselect"
	"cybex-gateway/worker/order"
	"cybex-gateway/worker/wx"
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
	go wx.HandleWorker(5)
	go jpworker.HandleWorker(5)
	jpselect.StartServer()
}
