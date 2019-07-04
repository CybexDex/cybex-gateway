package main

import (
	"os"

	"github.com/spf13/viper"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/server/jpselect"
	"cybex-gateway/utils/log"
	jpworker "cybex-gateway/worker/jpselect"
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
	go jpworker.HandleWorker(5)
	jpselect.StartServer()
	// user.StartServer()
}
