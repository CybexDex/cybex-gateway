package main

import (
	"fmt"

	"bitbucket.org/woyoutlz/bbb-gateway/utils/ecc"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
)

func main1() {
	config.LoadConfig("uat")
	s, _ := jp.DepositAddress("ETH")
	fmt.Println(s)
}
func main() {
	s := ecc.PriToPub("bf12996feeaa2977b6ca0d33a0e8bd2ccfc4844c6f8a7e6d15c099f8da4a255d")
	fmt.Println(s)
}
