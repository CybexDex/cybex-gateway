package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jporder"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/sha3"

	model "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/exevent"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/order"
	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
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
	Confirmations int64                  `json:"confirmations"`
	CreateAt      int64                  `json:"create_at"`
	UpdateAt      int64                  `json:"update_at"`
	Fee           string                 `json:"fee"`
	Data          map[string]interface{} `json:"data"`
	Hash          string                 `json:"hash"`
	ExtraData     string                 `json:"extraData"`
	Memo          string                 `json:"memo"`
	Timestamp     int64                  `json:"timestamp"`
	SendAgain     bool                   `json:"sendAgain"`
}

// Sig ...
type Sig struct {
	R string `json:"r"`
	S string `json:"s"`
	V int64  `json:"v"`
}

//OrderNotiRequest ...
type OrderNotiRequest struct {
	Result    *OrderNotiResult `json:"result"`
	Crypto    string           `json:"crypto"`
	Timestamp int64            `json:"timestamp"`
	Sig       *Sig             `json:"sig"`
}

//OrderNoti ...
func OrderNoti(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	utils.Infof("order noti request:\n %s", requestBody)

	request := OrderNotiRequest{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		utils.Errorf("Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	result := request.Result
	result.Timestamp = request.Timestamp

	if os.Getenv("env") != "dev" && os.Getenv("env") != "staging" {
		// verify sig
		pubKey := os.Getenv("pub_key")
		ok, err := verifySign(result, request.Sig, pubKey)
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
	assetRepo.GetByName(result.CoinType)
	exeventEntity := new(model.ExEvent)
	exeventEntity.AssetID = 1
	exeventEntity.Hash = result.Hash
	exeventEntity.JadepoolID = 1
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
		//todo:
		jporderEntity.AssetID = 1
		//todo:
		jporderEntity.JadepoolID = 1
		amount, condition, err := apd.NewFromString(result.Value)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.NewFromString error: %v, condition: %s", err, condition.String())
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		jporderEntity.Amount = amount
		err = jporderRepo.Create(jporderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create jporder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	} else {
		if jporderEntity.Status == result.State {
			// repeat request
			utils.Infof("repeat request, jadepool order id: %d", jporderEntity.JadepoolOrderID)
			tx.Commit()
			resp := utils.Message(true, "success")
			utils.Respond(w, resp)
			return
		}

		updateEntity := &model.JPOrder{}
		updateEntity.Status = result.State
		jporderEntity.Status = result.State
		err = jporderRepo.UpdateColumns(jporderEntity.ID, updateEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("Update jporder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	}

	// jporder has done, so create order
	if jporderEntity.Status == model.JPorderStatusDone {
		orderRepo := order.NewRepo(tx)
		orderEntity := new(model.Order)
		orderEntity.JPHash = jporderEntity.Hash
		orderEntity.Status = "INIT"
		orderEntity.Type = jporderEntity.Type
		orderEntity.JPUUHash = jporderEntity.UUHash
		orderEntity.AssetID = 1
		orderEntity.TotalAmount = jporderEntity.Amount
		orderEntity.Amount = jporderEntity.Amount
		fee, _, _ := apd.NewFromString("0")
		orderEntity.Fee = fee
		//todo:
		orderEntity.AppID = 1
		err = orderRepo.Create(orderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create order error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func verifySign(result *OrderNotiResult, sign *Sig, pubkey string) (bool, error) {
	buf, _ := json.Marshal(result)
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err := decoder.Decode(&obj)
	if err != nil {
		return false, err
	}

	pubKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return false, err
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}

	notiMsgStr := utils.BuildMsg(obj)
	sha3Hash := sha3.NewLegacyKeccak256()
	_, err = sha3Hash.Write([]byte(notiMsgStr))
	if err != nil {
		return false, err
	}
	msgBuf := sha3Hash.Sum(nil)

	// Decode hex-encoded serialized signature.
	decodedR, err := base64.StdEncoding.DecodeString(sign.R)
	if err != nil {
		return false, err
	}
	decodedS, err := base64.StdEncoding.DecodeString(sign.S)
	if err != nil {
		return false, err
	}
	signature := btcec.Signature{
		R: new(big.Int).SetBytes(decodedR),
		S: new(big.Int).SetBytes(decodedS),
	}
	return signature.Verify(msgBuf, pubKey), nil
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
