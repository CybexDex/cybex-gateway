package main

import (
	"os"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/server/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/server/user"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	go jp.StartServer()
	go user.StartServer()
	for {

	}
}
