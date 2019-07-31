package sass

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

	"github.com/stretchr/objx"
)

var api *API

// ParseSassNoti ...
func ParseSassNoti(noti []byte) (*OrderInfo, error) {
	if api == nil {
		InitAPI()
	}
	return api.ParseOrderNoti(noti)
}

// HandleWithdraw ...
func HandleWithdraw(result *OrderInfo) error {
	res, err := model.JPOrderFind(&model.JPOrder{
		BNOrderID: &result.ID,
		Current:   "jpsended",
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
		order.Hash = result.Hash
		order.UUHash = fmt.Sprintf("%s_%s_%d", result.Type, result.Hash, result.N)
		order.Confirmations = result.Confirmations
		order.CurrentState = strings.ToUpper(result.State)
		ordernow = order
	} else {
		err = fmt.Errorf("JPOrderFind Record lenth %d", lenRes)
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
func HandleDeposit(result *OrderInfo) (err error) {
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
func createJPOrderWithDeposit(result *OrderInfo) (*model.JPOrder, error) {
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
	asset := result.CoinName
	out, err := model.AssetsFrist(&model.Asset{
		JadeName: asset,
	})
	if err != nil {
		out, err = model.AssetsFrist(&model.Asset{
			Name: asset,
		})
		if err != nil {
			return nil, err
		}
	}
	fromArr := strings.Split(result.From, ",")
	theFrom := fromArr[0]
	jporder := &model.JPOrder{
		Asset:      out.Name,
		BlockChain: result.Type,
		BNOrderID:  &result.ID,
		CybUser:    user,
		OutAddr:    theFrom,
		From:       theFrom,
		To:         result.To,
		Memo:       result.Memo,

		Confirmations: result.Confirmations,
		// Resend:        result.SendAgain,

		Index:        result.N,
		Hash:         result.Hash,
		UUHash:       fmt.Sprintf("%s_%s_%d", result.Type, result.Hash, result.N),
		TotalAmount:  total,
		Type:         result.BizType,
		Status:       model.JPOrderStatusPending,
		Current:      "jp",
		CurrentState: strings.ToUpper(result.State),
	}
	jporder.CybAsset = out.CYBName
	// jporder.create
	// err = model.JPOrderCreate(jporder)
	return jporder, nil
}

// VerifyAddress ...
func VerifyAddress(asset string, address string) (res *types.VerifyRes, err error) {
	if api == nil {
		InitAPI()
	}
	re, err := api.VerifyAddress(asset, address)
	if err != nil {
		return nil, err
	}
	result := types.VerifyRes{}
	result.Asset = asset
	result.Address = address
	obj := objx.Map(re)
	result.Valid = obj.Get("valid").Bool()
	return &result, err
}

// InitAPI ...
func InitAPI() {
	api = &API{
		APIAddr:   viper.GetString("sassserver.host"),
		AppKey:    viper.GetString("sassserver.appKey"),
		AppSecret: viper.GetString("sassserver.appSecret"),
	}
}

// DepositAddress ...
func DepositAddress(coin string) (address *types.JPAddressResult, err error) {
	// 构造ecc部分
	if api == nil {
		InitAPI()
	}
	re, err := api.GetAdderess(coin)
	if err != nil {
		return nil, err
	}
	result := types.JPAddressResult{}
	result.Type = coin
	obj := objx.Map(re)
	result.Address = obj.Get("address").String()
	return &result, err
}

// Withdraw ...
func Withdraw(coin string, to string, value string, sequence uint) (address *types.JPOrderResult, err error) {
	// 构造data消息体
	if api == nil {
		InitAPI()
	}
	id := fmt.Sprintf("%s:%s:%s:%d", coin, to, value, sequence)
	re, err := api.Withdraw(id, coin, to, value, "")
	if err != nil {
		return nil, err
	}
	result := types.JPOrderResult{
		ID:            re.ID,
		CoinName:      re.CoinName,
		Txid:          re.Hash,
		State:         re.State,
		BizType:       re.BizType,
		Type:          re.Type,
		Fee:           re.Fee,
		Confirmations: re.Confirmations,
		N:             re.N,
		Memo:          re.Memo,
		Block:         re.Block,
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
