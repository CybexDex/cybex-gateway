package main

import (
	"flag"
	"fmt"

	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/cybsrv"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/spf13/viper"
)

var (
	version   string
	githash   string
	buildtime string
	branch    string
)

func main() {
	v := flag.Bool("v", false, "version")
	flag.Parse()
	if *v {
		fmt.Printf("version: %s_%s_%s, build time: %s\n", version, branch, githash, buildtime)
		return
	}

	// init config
	utils.InitConfig()
	logDir := viper.GetString("cybsrv.log_dir")
	logLevel := viper.GetString("cybsrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("version: %s_%s_%s, build time: %s", version, branch, githash, buildtime)
	// cybsrv.Test()
	go cybsrv.BlockRead()
	go cybsrv.HandleWorker()
	select {} // block forever
}
