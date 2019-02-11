package utils

import (
	"os"

	"github.com/spf13/viper"
)

// InitConfig ...
func InitConfig() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
