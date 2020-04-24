package main

import (
	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
	"cybex-gateway/worker/cybdotorder"
	"os"

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
	log.InitLog(logDir, logLevel, "[cybexdot]")

	model.INITFromViper()
	cybdotorder.InitNode()
	go cybdotorder.HandleWorker(5)
	go cybdotorder.BlockRead()
	select {}
	// fmt.Println(s)
}
