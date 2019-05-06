package main

import (
	"fmt"
	"os"

	"bitbucket.org/woyoutlz/bbb-gateway/config"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"github.com/shopspring/decimal"
)

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	model.INITFromViper()
	cyborder()
}
func cyborder() {
	amount, err := decimal.NewFromString("0.01")
	var db = model.GetDB()
	b := model.JPOrder{
		Asset:        "ETH",
		To:           "aaaaa",
		Amount:       amount,
		Current:      "cyborder",
		CurrentState: "INIT",
		CybUser:      "bbbusdtuser1-2",
		Type:         "DEPOSIT",
	}
	err = db.Create(&b).Error
	fmt.Println(err)
}
