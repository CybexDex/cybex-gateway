package jp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"

	"cybex-gateway/model"
	"cybex-gateway/types"
	"cybex-gateway/utils"
	"cybex-gateway/utils/ecc"
	"cybex-gateway/utils/log"
)

// HandleWithdraw ...
func HandleWithdraw(result types.JPOrderResult) error {
	res, err := model.JPOrderFind(&model.JPOrder{
		BNOrderID: &result.ID,
		Current:   "jpsended",
	})
	if err != nil {
		return utils.ErrorAdd(err, "HandleWithdraw")
	}
	lenRes := len(res)
	var ordernow *model.JPOrder
	if lenRes == 1 {
		order := res[0]
		if order.Status == model.JPOrderStatusDone || order.Status == model.JPOrderStatusFailed {
			// 如果已经是done或者fail,记录一条日志,返回错误
			log.Errorln("final order cannot change", order.ID)
			return fmt.Errorf("final order cannot change %d", order.ID)
		}
		orderSequence := order.ID*100 + order.BNRetry
		// 如果sequence不是当前的，直接返回
		if orderSequence != result.Sequence {
			return nil
		}
		order.Hash = result.Txid
		order.UUHash = fmt.Sprintf("%s_%s_%d", result.Type, result.Txid, result.N)
		order.Confirmations = result.Confirmations
		order.CurrentState = strings.ToUpper(result.State)
		order.Resend = result.SendAgain
		if order.Resend && order.CurrentState == model.JPOrderStatusFailed {
			// resend 的话增加retry,可以重发
			if order.BNRetry < 3 {
				order.BNRetry = order.BNRetry + 1
				order.Log("BnResend", fmt.Sprintf("bn:%+v,order:%+v", result, order))
				order.SetCurrent("jp", model.JPOrderStatusInit, fmt.Sprintf("BN resend,%d", order.BNRetry))
			} else {
				msg := fmt.Sprintf("gatewayID:%d,jadepoolID:%s", order.ID, *order.BNOrderID)
				model.WxSendTaskCreate("BN重发次数超过3次", msg)
			}
		}
		if order.Resend == false && order.CurrentState == model.JPOrderStatusFailed {
			msg := fmt.Sprintf("gatewayID:%d,jadepoolID:%s", order.ID, *order.BNOrderID)
			model.WxSendTaskCreate("BN failed,不resend", msg)
		}
		ordernow = order
	} else {
		err = fmt.Errorf("Record lenth %d", lenRes)
		log.Errorln(err, "resultID:", result.ID)
		return nil
	}
	if ordernow.CurrentState == model.JPOrderStatusDone {
		ordernow.SetCurrent("done", model.JPOrderStatusDone, "")
		ordernow.SetStatus(model.JPOrderStatusDone)
	}
	return ordernow.Save()
}

// HandleDeposit ...
func HandleDeposit(result types.JPOrderResult) (err error) {
	res, err := model.JPOrderFind(&model.JPOrder{
		BNOrderID: &result.ID,
		Current:   "jp",
	})
	if err != nil {
		return utils.ErrorAdd(err, "HandleDeposit")
	}
	lenRes := len(res)
	var ordernow *model.JPOrder
	if lenRes == 1 {
		order := res[0]
		if order.Status == model.JPOrderStatusDone || order.Status == model.JPOrderStatusFailed {
			// 如果已经是done或者fail,记录一条日志,返回错误
			log.Errorln("final order cannot change", order.ID)
			return fmt.Errorf("final order cannot change %d", order.ID)
		}
		order.Confirmations = result.Confirmations
		order.CurrentState = strings.ToUpper(result.State)
		ordernow = order
	} else if lenRes == 0 {
		// 创建订单,充值用户
		ordernow, err = createJPOrderWithDeposit(result)
		if err != nil {
			return utils.ErrorAdd(err, "HandleDeposit")
		}
	} else {
		return fmt.Errorf("Record lenth %d", lenRes)
	}
	if ordernow.CurrentState == model.JPOrderStatusDone {
		ordernow.SetCurrent("order", model.JPOrderStatusInit, "")
	}
	return ordernow.Save()
}
func createJPOrderWithDeposit(result types.JPOrderResult) (*model.JPOrder, error) {
	as, err := model.AddressFetch(&model.Address{
		Address: result.To,
	})
	if err != nil {
		return nil, err
	}
	if len(as) == 0 {
		return nil, fmt.Errorf("no_addrss_to_handle %s", result.To)
	}
	user := as[0].User
	total, err := decimal.NewFromString(result.Value)
	if err != nil {
		return nil, err
	}
	asset := result.Type
	if result.SubType != "" {
		asset = result.SubType
	}
	jporder := &model.JPOrder{
		Asset:      asset,
		BlockChain: result.Type,
		BNOrderID:  &result.ID,
		CybUser:    user,
		OutAddr:    result.From,
		From:       result.From,
		To:         result.To,
		Memo:       result.Memo,

		Confirmations: result.Confirmations,
		Resend:        result.SendAgain,

		Index:        result.N,
		Hash:         result.Txid,
		UUHash:       fmt.Sprintf("%s_%s_%d", result.Type, result.Txid, result.N),
		TotalAmount:  total,
		Type:         result.BizType,
		Status:       model.JPOrderStatusPending,
		Current:      "jp",
		CurrentState: strings.ToUpper(result.State),
	}
	// jporder.create
	// err = model.JPOrderCreate(jporder)
	return jporder, nil
}

// VerifyAddress ...
func VerifyAddress(asset string, address string) (res *types.VerifyRes, err error) {
	requestAddress := &types.JPAddressRequest{}
	requestAddress.Type = asset
	sendData, err := sendDataEcc(requestAddress)
	if err != nil {
		fmt.Println(err)
	}
	data := types.JPEvent{}
	log.Infoln(sendData.Data)
	err = bnResult("/api/v1/addresses/"+address+"/verify", sendData, &data)
	if err != nil {
		return nil, err
	}
	err = CheckComing(&data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	result := types.VerifyRes{}
	err = utils.ResultToStruct(data.Result, &result)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return &result, err
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
	log.Infoln(sendData.Data)
	err = bnResult("/api/v1/addresses/new", sendData, &data)
	if err != nil {
		return nil, err
	}
	err = CheckComing(&data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	result := types.JPAddressResult{}
	err = utils.ResultToStruct(data.Result, &result)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return &result, err
}

// Withdraw ...
func Withdraw(coin string, to string, value string, sequence uint) (address *types.JPOrderResult, err error) {
	// 构造data消息体
	requestObj := &types.JPWithdrawRequest{}
	requestObj.Type = coin
	requestObj.To = to
	requestObj.Value = value
	requestObj.Sequence = sequence
	// 构造ecc部分
	sendData, err := sendDataEcc(requestObj)
	if err != nil {
		fmt.Println(err)
	}
	data := types.JPEvent{}
	// fmt.Println(sendData.Data)
	err = bnResult("/api/v1/transactions", sendData, &data)
	if err != nil {
		return nil, err
	}
	err = CheckComing(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := types.JPOrderResult{}
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
