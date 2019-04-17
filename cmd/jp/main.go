package main

import (
	"os"

	"github.com/spf13/viper"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/server/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	jpworker "bitbucket.org/woyoutlz/bbb-gateway/worker/jp"
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
	jp.StartServer()
	// user.StartServer()
}
