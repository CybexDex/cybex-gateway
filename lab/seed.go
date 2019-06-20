package main

import (
	"cybex-gateway/config"
	"cybex-gateway/utils"
	"fmt"
)

func main() {
	config.LoadConfig("uat")
	s := utils.SeedString("seed__test2")
	fmt.Println(s)
}
