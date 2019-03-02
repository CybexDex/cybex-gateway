package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var inited = false

// InitConfig ...
func InitConfig() {
	if inited {
		return
	}
	env := os.Getenv("GATEWAY_ENV")
	if len(env) == 0 {
		env = "dev"
	}
	fmt.Println("env is ", env)
	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	inited = true
}
