package main

import (
	"cybex-gateway/worker/cybdotorder"
	"os"

	"github.com/spf13/viper"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/server/jpselect"
	"cybex-gateway/server/user"
	"cybex-gateway/utils/log"
	jpworker "cybex-gateway/worker/jpselect"
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

	cybdotorder.InitNode()
	go cybdotorder.HandleWorker(5)
	go cybdotorder.BlockRead()

	go user.StartServer()
	iswx := viper.GetBool("wx.enable")
	if iswx {
		go wx.HandleWorker(5)
	}
	go jpworker.HandleWorker(5)
	jpselect.StartServer()
}
