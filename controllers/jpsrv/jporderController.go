package jpsrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	model "coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/repository/address"
	"coding.net/bobxuyang/cy-gateway-BN/repository/app"
	"coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	"coding.net/bobxuyang/cy-gateway-BN/repository/exevent"
	"coding.net/bobxuyang/cy-gateway-BN/repository/jadepool"
	"coding.net/bobxuyang/cy-gateway-BN/repository/jporder"
	utils "coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

//OrderNotiResult ...
type OrderNotiResult struct {
	ID            string                 `json:"id"`
	Sequence      uint                   `json:"sequence"`
	State         string                 `json:"state"`
	BizType       string                 `json:"bizType"`
	CoinType      string                 `json:"coinType"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
	Value         string                 `json:"value"`
	Confirmations int                    `json:"confirmations"`
	CreateAt      int64                  `json:"create_at"`
	UpdateAt      int64                  `json:"update_at"`
	Fee           string                 `json:"fee"`
	Data          map[string]interface{} `json:"data"`
	Hash          string                 `json:"hash"`
	ExtraData     string                 `json:"extraData"`
	Memo          string                 `json:"memo"`
	Timestamp     int64                  `json:"timestamp"`
	SendAgain     bool                   `json:"sendAgain"`
	Namespace     string                 `json:"namespace,omitempty"`
	Sid           string                 `json:"sid,omitempty"`
}

//JPComeData ...
type JPComeData struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Status     int                    `json:"status"`
	Result     map[string]interface{} `json:"result"`
	Crypto     string                 `json:"crypto"`
	Timestamp  int64                  `json:"timestamp"`
	Sig        utils.ECCSig           `json:"sig"`
	JadepoolID int                    `json:"jadepoolID"`
}

//NotiOrder ...
func NotiOrder(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	utils.Infof("order noti request:\n %s", requestBody)

	request := JPComeData{}
	decoder := json.NewDecoder(bytes.NewReader(requestBody))
	decoder.UseNumber()
	err = decoder.Decode(&request)
	if err != nil {
		utils.Errorf("Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	request.Result["timestamp"] = request.Timestamp

	jadepoolID := uint(request.JadepoolID)
	// if jadepoolID not exist, use default 1
	if jadepoolID == 0 {
		jadepoolID = 1
	}
	if os.Getenv("env") != "dev" {
		// verify sig
		jadepoolRepo := jadepool.NewRepo(model.GetDB())
		jadepool, err := jadepoolRepo.GetByID(jadepoolID)
		if err != nil {
			utils.Errorf("Unmarshal error: %v", err)
			utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}
		pubKey := jadepool.EccPubKey
		ok, err := utils.VerifyECCSign(request.Result, &request.Sig, pubKey)
		if err != nil {
			utils.Errorf("verifySign error: %v", err)
			utils.Respond(w, utils.Message(false, "Sign error"), http.StatusForbidden)
			return
		}
		if !ok {
			utils.Errorf("verify result: %v", ok)
			utils.Respond(w, utils.Message(false, "Sign error"), http.StatusForbidden)
			return
		}
	}

	resultBytes, err := json.Marshal(request.Result)
	if err != nil {
		utils.Errorf("Marshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	result := &OrderNotiResult{}
	err = json.Unmarshal(resultBytes, result)
	if err != nil {
		utils.Errorf("Marshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	result.State = strings.ToUpper(result.State)

	// begin transaction
	tx := model.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			utils.Errorf("%v, stack: %s", r, debug.Stack())
			tx.Rollback()
		}
	}()

	// save exevent
	assetRepo := asset.NewRepo(tx)
	asset, err := assetRepo.GetByName(result.CoinType)
	if err != nil {
		utils.Errorf("assetRepo.GetByName error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	exeventEntity := new(model.ExEvent)
	exeventEntity.AssetID = asset.ID
	exeventEntity.Hash = result.Hash
	exeventEntity.JadepoolID = jadepoolID
	exeventEntity.Status = result.State
	exeventEntity.Log = string(requestBody)
	exeventRepo := exevent.NewRepo(tx)
	err = exeventRepo.Create(exeventEntity)
	if err != nil {
		tx.Rollback()
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	// find jporder and save/update jporder
	jpOrderID, err := strconv.Atoi(result.ID)
	if err != nil {
		tx.Rollback()
		utils.Errorf("atoi error: %v, id: %s", err, result.ID)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	jporderRepo := jporder.NewRepo(tx)
	jporderEntity, err := jporderRepo.GetByJPID(uint(jpOrderID))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		utils.Errorf("get jporder error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	if jporderEntity == nil {
		addressRepo := address.NewRepo(tx)
		adddresses, err := addressRepo.FetchWith(&model.Address{
			Address: result.To,
		})
		if err != nil {
			tx.Rollback()
			utils.Errorf("addressRepo.FetchWith error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		if len(adddresses) == 0 {
			tx.Rollback()
			utils.Errorf("addressRepo.FetchWith result: %v", adddresses)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		appID := adddresses[0].AppID

		// parse tx index from result
		n := parseIndexFromResult(result)
		jporderEntity = new(model.JPOrder)
		jporderEntity.From = result.From
		jporderEntity.To = result.To
		jporderEntity.Hash = result.Hash
		if n > 0 {
			jporderEntity.UUHash = fmt.Sprintf("%s:%s:%d", result.CoinType, result.Hash, n)
		} else {
			jporderEntity.UUHash = fmt.Sprintf("%s:%s", result.CoinType, result.Hash)
		}

		jporderEntity.Index = n
		jporderEntity.JadepoolOrderID = uint(jpOrderID)
		jporderEntity.Status = result.State
		if jporderEntity.Status != model.JPOrderStatusDone && jporderEntity.Status != model.JPOrderStatusFailed {
			jporderEntity.Status = model.JPOrderStatusPending
		}
		jporderEntity.Type = result.BizType
		jporderEntity.AssetID = asset.ID
		jporderEntity.AppID = appID
		jporderEntity.JadepoolID = jadepoolID
		totalAmount, condition, err := apd.NewFromString(result.Value)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.NewFromString error: %v, condition: %s", err, condition.String())
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		jporderEntity.TotalAmount = totalAmount
		if jporderEntity.Type == model.JPOrderTypeDeposit && jporderEntity.TotalAmount.Cmp(asset.DepositFee) < 0 {
			tx.Rollback()
			utils.Errorf("total is smaller than fee: %v < %v", jporderEntity.TotalAmount, asset.DepositFee)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		jporderEntity.Amount = apd.New(0, 0)
		condition, err = apd.BaseContext.Sub(jporderEntity.Amount, jporderEntity.TotalAmount, asset.DepositFee)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.BaseContext.Sub error: %v, condition: %v", err, condition)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}

		jporderEntity.Fee = asset.DepositFee
		jporderEntity.Confirmations = result.Confirmations
		jporderEntity.Resend = result.SendAgain

		err = jporderRepo.Create(jporderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create jporder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}

		jporderEntity.EnterHook = true
		err = jporderEntity.AfterSaveHook(tx)
		if err != nil {
			tx.Rollback()
			utils.Errorf("call jporder after-save hook error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	} else {
		if jporderEntity.Status == model.JPOrderStatusDone ||
			jporderEntity.Status == model.JPOrderStatusFailed {
			// repeat request
			utils.Infof("repeat request, jadepool order id: %d", jporderEntity.JadepoolOrderID)
			tx.Commit()
			resp := utils.Message(true, "success")
			utils.Respond(w, resp)
			return
		}

		updateEntity := &model.JPOrder{}
		if result.State != model.JPOrderStatusInit {
			updateEntity.Status = result.State
		}
		if updateEntity.Status != model.JPOrderStatusDone && updateEntity.Status != model.JPOrderStatusFailed {
			updateEntity.Status = model.JPOrderStatusPending
		}

		updateEntity.Confirmations = result.Confirmations
		updateEntity.Resend = result.SendAgain
		if jporderEntity.Hash == "" {
			jporderEntity.Hash = result.Hash
		}
		jporderEntity.Status = updateEntity.Status
		err = jporderRepo.UpdateColumns(jporderEntity.ID, updateEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("Update jporder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}

		jporderEntity.EnterHook = true
		err = jporderEntity.AfterSaveHook(tx)
		if err != nil {
			tx.Rollback()
			utils.Errorf("call jporder after-save hook error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

///////////////////////////////////////////////send transaction///////////////////////////////////////////////

// JPSendRequest ...
type JPSendRequest struct {
	ID uint `json:"id"`
}

// JPTransaction ...
type JPTransaction struct {
	Sequence  uint   `json:"sequence"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	To        string `json:"to"`
	Timestamp int64  `json:"timestamp"`
	Callback  string `json:"callback"`
	ExtraData string `json:"extraData"`
}

// JPSendData ...
type JPSendData struct {
	Crypto    string        `json:"crypto"`
	Hash      string        `json:"hash"`
	Encode    string        `json:"encode"`
	AppID     string        `json:"appid"`
	Timestamp int64         `json:"timestamp"`
	Sig       *utils.ECCSig `json:"sig"`
	Data      interface{}   `json:"data"`
}

// SendOrder ...
func SendOrder(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "bad request error"), http.StatusBadRequest)
		return
	}
	reqData := JPSendRequest{}
	err = json.Unmarshal(bodyBytes, &reqData)
	if err != nil {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "bad request error"), http.StatusBadRequest)
		return
	}

	result, err := DoSendOrder(reqData.ID)
	if err != nil {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	respData := utils.Message(true, "success", result)
	utils.Respond(w, respData)
}

// DoSendOrder ...
func DoSendOrder(id uint) (map[string]interface{}, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}
	jporderRepo := jporder.NewRepo(model.GetDB())
	jporder, err := jporderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	if jporder.AssetID == 0 || len(jporder.To) == 0 {
		return nil, fmt.Errorf("error: %v", err)
	}

	assetRepo := asset.NewRepo(model.GetDB())
	asset, err := assetRepo.GetByID(jporder.AssetID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	timestamp := time.Now().Unix() * 1000
	jptransaction := &JPTransaction{}
	jptransaction.Sequence = jporder.ID
	jptransaction.Type = asset.Name
	jptransaction.Value = jporder.Amount.String()
	jptransaction.To = jporder.To
	jptransaction.Timestamp = timestamp
	jptransaction.Callback = viper.GetString("jpsrv.self_addr") + "/api/order/noti"

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = viper.GetString("jpsrv.jadepool_appid")
	sendData.Data = &jptransaction

	appID := jporder.AppID
	appRepo := app.NewRepo(model.GetDB())
	app, err := appRepo.GetByID(appID)
	if err != nil {
		return nil, fmt.Errorf("appRepo.GetByID error: %v", err)
	}
	jadepoolID := app.JadepoolID
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(jadepoolID)
	if err != nil {
		return nil, fmt.Errorf("jadepoolRepo.GetByID error: %v, id: %d", err, jadepoolID)
	}

	prikey := viper.GetString("jpsrv.pri_key")
	sig, err := utils.SignECCData(prikey, sendData.Data)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	url := jadepoolAddr + "/api/v1/transactions/"

	orderResp := JPComeData{}
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("post error: %v", err)
	}
	_bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	err = json.Unmarshal(_bodyBytes, &orderResp)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v, body: %s", err, string(_bodyBytes))
	}

	if orderResp.Code != 0 || orderResp.Status != 0 || orderResp.Result == nil {
		return nil, fmt.Errorf("response: %#v", orderResp)
	}

	if os.Getenv("env") != "dev" {
		// verify sig
		orderResp.Result["timestamp"] = orderResp.Timestamp
		pubKey := jadepool.EccPubKey
		ok, err := utils.VerifyECCSign(orderResp.Result, &orderResp.Sig, pubKey)
		if err != nil {
			return nil, fmt.Errorf("verifySign error: %v, data: %#v", err, orderResp)
		}
		if !ok {
			return nil, fmt.Errorf("verify result: %v, data: %#v", ok, orderResp)
		}
	}

	return orderResp.Result, nil
}

func parseIndexFromResult(result *OrderNotiResult) int {
	n := 0
	if result.CoinType != "BTC" && result.CoinType != "QTUM" {
		return n
	}

	toes := result.Data["to"]
	if toes == nil {
		return n
	}

	switch toes.(type) {
	case []interface{}:
		toesMap := toes.([]interface{})
		if len(toesMap) <= 1 {
			break
		}
		for _, v := range toesMap {
			val := v.(map[string]interface{})
			if val["address"].(string) != result.To {
				continue
			}
			switch val["n"].(type) {
			case float64:
				n = int(val["n"].(float64))
			case int64:
				n = int(val["n"].(int64))
			case int32:
				n = int(val["n"].(int32))
			}
			break
		}
	default:
		utils.Errorf("data.to type  is: %v", reflect.TypeOf(toes))
	}
	return n
}
