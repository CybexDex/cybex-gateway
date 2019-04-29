package main

import (
	"fmt"
	"os"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	model.INITFromViper()
	black()
	asset()
}
func asset() {
	var db = model.GetDB()
	b := model.Asset{
		Name:         "ETH",
		Blockchain:   "ETH",
		CYBName:      "JADE.ETH",
		Confirmation: "20",

		SmartContract:  "",
		GatewayAccount: "bbbusdtin1",
		WithdrawPrefix: "withdraw:CybexGatewayDev",

		DepositSwitch:  true,
		WithdrawSwitch: true,

		// MinDeposit  :"0",
		// MinWithdraw :"0",
		// WithdrawFee :"0",
		// DepositFee :"0",

		// Precision	:"",
		// ImgURL :"",
		HashLink: "https://etherscan.io/tx/%s",
	}
	err := db.Create(&b).Error
	fmt.Println(err)
}
func black() {
	var db = model.GetDB()
	b := model.Black{
		Blockchain: "CYB",
		Address:    "yangyu3",
	}
	err := db.Create(&b).Error
	fmt.Println(err)
	b = model.Black{
		Blockchain: "ETH",
		Address:    "0xb8a51ef04e0f4ca102eff710f534c2b9509ca1e3",
	}
	err = db.Create(&b).Error
	fmt.Println(err)
}
