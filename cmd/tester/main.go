package main

import (
	"encoding/json"
	"fmt"

	"cybex-gateway/model"
	"cybex-gateway/types"
	"cybex-gateway/utils/ecc"

	"cybex-gateway/config"
	"cybex-gateway/controller/jp"
	"cybex-gateway/controller/bbb"
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
func mainNewPriPub() {
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
func main4() {
	config.LoadConfig("dev")
	sign()
}
func bbbInfo(){
	fmt.Println("bbbInfo")
	x,err := bbb.Info()
	fmt.Println(x.BlockNum,x.BlockID,err)
}
func bbbSend(){
	x := `{"signatures":["1f3a31b04b3ecd1840181767bd4303a84c786a02d35413730c3dadfdb5b6eeedf93955a1bafb05757e7eb98af8224bfc27452287d71374302ee57237ae858cffc0"],"ref_block_num":62161,"ref_block_prefix":1736324565,"expiration":"2020-02-10T02:16:34","operations":[[0,{"from":"1.2.68142","to":"1.2.38759","amount":{"amount":100000,"asset_id":"1.3.27"},"memo":{"from":"CYB5onfmMrkPAoUVmjVQ9Lk1fL9zm1f6hdHzmCvr5Ust5qMRBUQAu","to":"CYB53cm9QCfhsHUF1kXcY1YiZfqc3Mq3qCtYusW3xhpJyJsdGvbYu","nonce":5577006791947779410,"message":"cfcb4aacfd9521c1323c623bb1ebeba8159429833f00bc9d6088ce2079e41ba8f29b9227deac9b352372e3f2f4da38f8"},"extensions":[],"fee":{"amount":1067,"asset_id":"1.3.0"}}],[0,{"from":"1.2.38759","to":"1.2.68143","amount":{"amount":100000,"asset_id":"1.3.1145"},"memo":{"from":"CYB53cm9QCfhsHUF1kXcY1YiZfqc3Mq3qCtYusW3xhpJyJsdGvbYu","to":"CYB5Wxd3T7nSt3FQnb3LMAPxaH8SQFkg447Xwj6DT8w8dXf4ivrac","nonce":8674665223082153551,"message":"9c771ecb7ab693de570d73c7c59e0571c69b36b1f531dd715f96e15f7193c04f2ac1b318456193ba98d0793467bfe4f9"},"extensions":[],"fee":{"amount":1067,"asset_id":"1.3.0"}}]],"extensions":[]}`
	err := bbb.BroadcastTransaction(x)
	fmt.Println(err)
}
func main(){
	config.LoadConfig("bbb")
	// bbbInfo()
	bbbSend()
}