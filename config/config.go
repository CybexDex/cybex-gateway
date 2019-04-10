package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig ...
func LoadConfig(env string) {
	envlocal := env
	fmt.Println("try to find", envlocal)
	viper.SetConfigName(envlocal)
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	viper.Get("")
}

//SeedString ...
func SeedString(path string) string {
	return ""
}
