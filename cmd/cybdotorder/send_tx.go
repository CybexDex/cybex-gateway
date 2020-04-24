package main

import (
	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
	"cybex-gateway/worker/cybdotorder"
	"os"

	"github.com/centrifuge/go-substrate-rpc-client/types"

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

	hash, err := types.NewHashFromHexString("0xf4169148358073831eacb40822ccfa8a7754c8fd8e5283be0dc98db8e86181ec")
	if err != nil {
		log.Errorln(err)
	}
	extrinsic, err := cybdotorder.CreateTransfer("5ELsycSf2vVGN4whjWA176J2TWzJY52PDGxTC89APeuTcbvB", 1000, hash, "xxx")
	if err != nil {
		log.Errorln(err)
	}
	txHash, err := cybdotorder.SignAndSendTransfer(extrinsic, "evil person wolf innocent blast panther someone police rhythm target sheriff toddler")
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln(txHash)
}
