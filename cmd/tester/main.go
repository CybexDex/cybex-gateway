package main

import (
	"encoding/json"
	"fmt"

	"cybex-gateway/model"
	"cybex-gateway/types"
	"cybex-gateway/utils/ecc"

	"cybex-gateway/config"
	"cybex-gateway/controller/jp"
	"cybex-gateway/controller/user"
)

func main1() {
	config.LoadConfig("dev")
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
func main4() {
	pri, pub := ecc.NewPriPub()
	fmt.Println(pri, pub)
}
func main5() {
	ecc.TestECCSign()
}
func sign() {
	requestObj := &types.JPWithdrawRequest{}
	// requestObj.Type = coin
	requestObj.To = "0x3ae306d3fe3584ec90a765db587815b6d990ce4a"
	requestObj.Value = "0.15"
	requestObj.Sequence = 57000
	// 构造ecc部分
	timestamp := int64(1562745578000)

	requestObj.Timestamp = timestamp
	// 构造ecc部分
	sendData, _ := jp.SendDataEcc(requestObj, timestamp)
	b, err := ecc.VerifyECCSign(sendData.Data, sendData.Sig, "02fe13e7db6d8098ac3c744291ea9108632a214f81c17e974c764d75f4ab7aa2b9")
	fmt.Println(b, err)
	bs, _ := json.Marshal(sendData)
	strsend := string(bs)
	fmt.Println("send jp json", strsend)
}
func main() {
	config.LoadConfig("dev")
	sign()
}
