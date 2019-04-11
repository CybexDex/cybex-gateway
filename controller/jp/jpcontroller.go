package jp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/ecc"
)

// HandleWithdraw ...
func HandleWithdraw(result types.JPOrderResult) error {
	// 事务
	// 更新提现订单
	// jporderID := result.ID
	// 寻找提现订单
	// jporder.update
	// 告知record
	// record.afterJPWithdraw
	// 记录Done事件
	// 或者抛出错误
	return nil
}

// HandleDeposit ...
func HandleDeposit(result types.JPOrderResult) error {
	// 事务
	// 新建充值订单
	// 充值地址
	// jporder := &m.JPOrder{
	// 	Asset:      result.Type,
	// 	BlockChain: result.SubType,
	// }
	// jporder.create
	// 或者更新充值订单
	// jporder.update
	// 告知record
	// record.afterJPDeposit
	// 记录Done事件
	// 或者抛出错误
	fmt.Println("deposit", result)
	return nil
}

// DepositAddress ...
func DepositAddress(coin string) (address *types.JPAddressResult, err error) {
	// 构造data消息体
	requestAddress := &types.JPAddressRequest{}
	requestAddress.Type = coin
	// 构造ecc部分
	sendData, err := sendDataEcc(requestAddress)
	if err != nil {
		fmt.Println(err)
	}
	data := types.JPEvent{}
	fmt.Println(sendData.Data)
	err = bnResult("/api/v1/addresses/new", sendData, &data)
	if err != nil {
		return nil, err
	}
	err = CheckComing(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := types.JPAddressResult{}
	err = utils.ResultToStruct(data.Result, &result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &result, err
}

//CheckComing ...
func CheckComing(data *types.JPEvent) (err error) {
	// checkEcc
	if data.Code != 0 || data.Status != 0 || data.Result == nil {
		return fmt.Errorf("BN request failed, data: %#v", data)
	}
	if viper.GetBool("jpserver.ecc") == true {
		// verify sig
		data.Result["timestamp"] = data.Timestamp
		pubKey := utils.SeedString(viper.GetString("jpserver.eccPub"))
		ok, err := ecc.VerifyECCSign(data.Result, data.Sig, pubKey)
		if err != nil {
			return fmt.Errorf("verifySign error: %v, data: %#v", err, data)
		}
		if !ok {
			return fmt.Errorf("verify result: %v, data: %#v", ok, data)
		}
	}
	return nil
}
func bnResult(urlPath string, sendData *types.JPSendData, v interface{}) (err error) {
	bs, _ := json.Marshal(sendData)
	jadepoolAddr := viper.GetString("jpserver.bnhost")
	url := jadepoolAddr + urlPath
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("ReadAll error: %v", err)
	}
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	return nil
}
func sendDataEcc(data interface{}) (sendData *types.JPSendData, err error) {
	sendData = &types.JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	timestamp := time.Now().Unix() * 1000
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = viper.GetString("jpserver.appid")
	sendData.Data = data
	//签名
	priKey := utils.SeedString(viper.GetString("jpserver.eccPri"))
	sig, err := ecc.SignECCData(priKey, data)
	if err != nil {
		return nil, fmt.Errorf("SignECCData error: %v", err)
	}
	sendData.Sig = sig
	return sendData, nil
}
