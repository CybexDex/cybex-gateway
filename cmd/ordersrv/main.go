package main

import (
	"flag"
	"fmt"

	"git.coding.net/bobxuyang/cy-gateway-BN/controllers/ordersrv"
	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	model "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	utils.InitConfig()
	logDir := viper.GetString("ordersrv.log_dir")
	logLevel := viper.GetString("ordersrv.log_level")
	// init logger
	utils.InitLog(logDir, logLevel)
	utils.Infof("version: %s_%s_%s, build time: %s", version, branch, githash, buildtime)

	// init db
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	model.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	rep.Init()

	ordersrv.HandleWorker()
}
