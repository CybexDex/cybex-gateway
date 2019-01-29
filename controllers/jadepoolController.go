package controllers

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

	"git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jadepool"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jporder"

	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"

	model "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/exevent"
	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

const (
	defaultJadepoolName = "Jadepool-001"
)

var (
	defaultJadepool *model.Jadepool
)

//OrderNotiResult ...
type OrderNotiResult struct {
	ID            string                 `json:"id"`
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
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Status    int                    `json:"status"`
	Result    map[string]interface{} `json:"result"`
	Crypto    string                 `json:"crypto"`
	Timestamp int64                  `json:"timestamp"`
	Sig       utils.ECCSig           `json:"sig"`
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

	if os.Getenv("env") != "dev" {
		// verify sig
		pubKey := os.Getenv("pub_key")
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
	jadepoolRepo := jadepool.NewRepo(tx)
	asset, err := assetRepo.GetByName(result.CoinType)
	if err != nil {
		utils.Errorf("assetRepo.GetByName error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	if defaultJadepool == nil {
		jadepool, err := jadepoolRepo.GetByName(defaultJadepoolName)
		if err != nil {
			utils.Errorf("assetRepo.GetByName error: %v", err)
			utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}
		defaultJadepool = jadepool
	}

	exeventEntity := new(model.ExEvent)
	exeventEntity.AssetID = asset.ID
	exeventEntity.Hash = result.Hash
	exeventEntity.JadepoolID = defaultJadepool.ID
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

	appID := uint(1)
	/*addressRepo := address.NewRepo(tx)
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
	appID = adddresses[0].AppID*/

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
		jporderEntity.Type = result.BizType
		jporderEntity.AssetID = asset.ID
		jporderEntity.AppID = appID
		jporderEntity.JadepoolID = defaultJadepool.ID
		totalAmount, condition, err := apd.NewFromString(result.Value)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.NewFromString error: %v, condition: %s", err, condition.String())
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		jporderEntity.TotalAmount = totalAmount
		if jporderEntity.TotalAmount.Cmp(asset.DepositFee) < 0 {
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
	} else {
		if jporderEntity.Status == result.State &&
			(jporderEntity.Status == model.JPOrderStatusDone ||
				jporderEntity.Status == model.JPOrderStatusFailed) {
			// repeat request
			utils.Infof("repeat request, jadepool order id: %d", jporderEntity.JadepoolOrderID)
			tx.Commit()
			resp := utils.Message(true, "success")
			utils.Respond(w, resp)
			return
		}

		updateEntity := &model.JPOrder{}
		updateEntity.Status = result.State
		updateEntity.Confirmations = result.Confirmations
		updateEntity.Resend = result.SendAgain
		jporderEntity.Status = result.State
		err = jporderRepo.UpdateColumns(jporderEntity.ID, updateEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("Update jporder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

///////////////////////////////////////////////send transaction///////////////////////////////////////////////

// JPTransaction ...
type JPTransaction struct {
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
	jptransaction := JPTransaction{}
	err = json.Unmarshal(bodyBytes, &jptransaction)
	if err != nil {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "bad request error"), http.StatusBadRequest)
		return
	}
	if len(jptransaction.Type) == 0 || len(jptransaction.Value) == 0 || len(jptransaction.To) == 0 {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "bad request error"), http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Unix() * 1000
	jptransaction.Timestamp = timestamp
	jptransaction.Callback = os.Getenv("self_url") + "/api/order/noti"

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = "app"
	sendData.Data = &jptransaction

	prikey := os.Getenv("pri_key")
	sig, err := utils.SignECCData(prikey, sendData.Data)
	if err != nil {
		utils.Errorf("create jporder error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	jadepoolURL := os.Getenv("jadepool_url")
	url := jadepoolURL + "/api/v1/transactions/"

	orderResp := JPComeData{}
	for i := 0; i < 3; i++ {
		resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
		if err != nil {
			utils.Errorf("post error: %v", err)
			continue
		}
		_bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			utils.Errorf("ReadAll error: %v", err)
			continue
		}

		err = json.Unmarshal(_bodyBytes, &orderResp)
		if err != nil {
			utils.Errorf("Unmarshal error: %v, body: %s", err, string(_bodyBytes))
			continue
		}
		break
	}

	if orderResp.Code != 0 || orderResp.Status != 0 || orderResp.Result == nil {
		utils.Errorln("not found address")
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", orderResp.Result)
	utils.Respond(w, resp)
}

/////////////////////////////////////////////address/////////////////////////////////////////

// JPAddressRequest ...
type JPAddressRequest struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Callback  string `json:"callback,omitempty"`
}

// JPAddressResult ...
type JPAddressResult struct {
	Type      string `json:"type"`
	SubType   string `json:"subType"`
	Address   string `json:"address"`
	Namespace string `json:"namespace,omitempty"`
	Sid       string `json:"sid,omitempty"`
}

//GetNewAddress ...
func GetNewAddress(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	coinType := query.Get("type")
	if len(coinType) == 0 {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	coinType = strings.ToUpper(coinType)

	timestamp := time.Now().Unix() * 1000
	requestAddress := JPAddressRequest{}
	requestAddress.Timestamp = timestamp
	requestAddress.Type = coinType
	requestAddress.Callback = os.Getenv("self_url") + "/api/order/noti"

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = "app"
	sendData.Data = &requestAddress

	prikey := os.Getenv("pri_key")
	sig, err := utils.SignECCData(prikey, sendData.Data)
	if err != nil {
		utils.Errorf("create jporder error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	jadepoolURL := os.Getenv("jadepool_url")
	url := jadepoolURL + "/api/v1/addresses/new"

	data := JPComeData{}
	for i := 0; i < 3; i++ {
		resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
		if err != nil {
			utils.Errorf("post error: %v", err)
			continue
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			utils.Errorf("ReadAll error: %v", err)
			continue
		}

		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			utils.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
			continue
		}
		break
	}
	if data.Code != 0 || data.Status != 0 || data.Result == nil || data.Result["address"] == nil {
		utils.Errorf("not found address, data: %#v", data)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	if os.Getenv("env") != "dev" {
		// verify sig
		pubKey := os.Getenv("pub_key")
		ok, err := utils.VerifyECCSign(data.Result, &data.Sig, pubKey)
		if err != nil {
			utils.Errorf("verifySign error: %v, data: %#v", err, data)
			utils.Respond(w, utils.Message(false, "Sign error"), http.StatusForbidden)
			return
		}
		if !ok {
			utils.Errorf("verify result: %v, data: %#v", ok, data)
			utils.Respond(w, utils.Message(false, "Sign error"), http.StatusForbidden)
			return
		}
	}

	resp := utils.Message(true, "success", data.Result)
	utils.Respond(w, resp)
}

func parseIndexFromResult(result *OrderNotiResult) int {
	n := 0
	if result.CoinType != "BTC" && result.CoinType != "QTUM" {
		return n
	}

	toes := result.Data["to"]
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
