package main

import (
	"os"

	"cybex-gateway/config"
	"cybex-gateway/model"
	"cybex-gateway/worker/order"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	model.INITFromViper()
	order.HandleWorker(5)
	// fmt.Println(s)
}
