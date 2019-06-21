package main

import (
	"cybex-gateway/utils"
	"fmt"
	"time"

	apim "github.com/CybexDex/cybex-go/api"
)

func main() {
	timestamp := time.Now().Unix()

	fmt.Println(timestamp)
	user := "a-new-world"
	password := "P5K5BmnCnFkX8gQ9viAXxVK5PUz3CiKZowXCD2m9MZWygo5GrWmn"
	pubkey := "CYB7b4XjdtiHMZ6xmkLJDXLxQ4Zbcwr6MQU4zctkMJzDf1DH4TXSM"
	to := fmt.Sprintf("%d%s", timestamp, user)

	str, err := utils.CybSign(user, pubkey, password, to)
	fmt.Println(str, err)
	api := apim.New("ws://18.136.140.223:38090/", "")
	if err := api.Connect(); err != nil {
		panic(err)
	}
	x, err := api.VerifySign(user, to, str)
	fmt.Println(x, err)
}

// func verify() {
// 	VerifySign("xiao01", toSign, sig)
// }
