package main

import (
	"os"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/worker/cyborder"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	model.INITFromViper()
	cyborder.InitNode()
	cyborder.InitAsset()
	// go cyborder.HandleWorker(5)
	cyborder.Test()
	// select {}
	// fmt.Println(s)
}
