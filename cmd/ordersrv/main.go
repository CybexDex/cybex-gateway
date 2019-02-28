package main

import (
	"flag"
	"fmt"

	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/ordersrv"
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

	utils.InitConfig()
	logDir := viper.GetString("ordersrv.log_dir")
	logLevel := viper.GetString("ordersrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("version: %s_%s_%s, build time: %s", version, branch, githash, buildtime)
	ordersrv.HandleWorker()
}
