package main

import (
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/ordersrv"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/spf13/viper"
)

var (
	githash   string
	buildtime string
	branch    string
)

func main() {
	utils.InitConfig()
	logDir := viper.GetString("ordersrv.log_dir")
	logLevel := viper.GetString("ordersrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)
	ordersrv.HandleWorker()
}
