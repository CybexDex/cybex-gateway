package main

import (
	"fmt"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/ecc"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/controller/user"
)

func main1() {
	config.LoadConfig("uat")
	s, _ := jp.DepositAddress("ETH")
	fmt.Println(s)
}
func main2() {
	s := ecc.PriToPub("bf12996feeaa2977b6ca0d33a0e8bd2ccfc4844c6f8a7e6d15c099f8da4a255d")
	fmt.Println(s)
}
func main3() {
	config.LoadConfig("uat")
	model.INITFromViper()
	s, err := user.GetAddress("yangyu4", "ETH1")
	fmt.Println(s, err)
}
func main() {
	ecc.TestECCSign()
}
