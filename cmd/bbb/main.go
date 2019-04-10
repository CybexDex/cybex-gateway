package main

import (
	"os"

	"github.com/spf13/viper"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/server/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/server/user"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
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
	go jp.StartServer()
	go user.StartServer()
	for {

	}
}
