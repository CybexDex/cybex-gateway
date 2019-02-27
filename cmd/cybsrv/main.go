package main

import (
	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/cybsrv"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/spf13/viper"
)

var (
	githash   string
	buildtime string
	branch    string
)

func main() {
	// init config
	utils.InitConfig()
	logDir := viper.GetString("cybsrv.log_dir")
	logLevel := viper.GetString("cybsrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("build info: %s_%s_%s", buildtime, branch, githash)
	// cybsrv.Test()
	go cybsrv.BlockRead()
	go cybsrv.HandleWorker()
	select {} // block forever
}
