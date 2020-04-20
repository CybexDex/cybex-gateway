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

	extrinsic, err := cybdotorder.MakeTransfer("5HQXLYqiiisunFjNvc164QSajMyytzgaqJJSd7bCx82bxi6W", 1000)
	if err != nil {
		log.Errorln(err)
	}
	txHash, err := cybdotorder.SignAndSendTransfer(extrinsic, "staff mammal myself patrol notice neglect pass shine scale cliff nominee popular")
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln(txHash)
}
