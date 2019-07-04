package main

import (
	"fmt"

	"cybex-gateway/model"
	"cybex-gateway/utils/ecc"

	"cybex-gateway/config"
	"cybex-gateway/controller/jp"
	"cybex-gateway/controller/user"
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
func main3() {
	config.LoadConfig("uat")
	model.INITFromViper()
	s, err := user.GetAddress("yangyu4", "ETH1")
	fmt.Println(s, err)
}
func main4() {
	pri, pub := ecc.NewPriPub()
	fmt.Println(pri, pub)
}
func main5() {
	ecc.TestECCSign()
}
