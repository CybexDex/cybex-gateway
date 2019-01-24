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

	"github.com/btcsuite/btcd/btcec"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/sha3"

	model "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/exevent"
	exorder "git.coding.net/bobxuyang/cy-gateway-BN/repository/exorder"
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

	// find exorder and save/update exorder
	jpOrderID, err := strconv.Atoi(result.ID)
	if err != nil {
		tx.Rollback()
		utils.Errorf("atoi error: %v, id: %s", err, result.ID)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	exorderRepo := exorder.NewRepo(tx)
	exorderEntity, err := exorderRepo.GetByJPID(uint(jpOrderID))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		utils.Errorf("get exorder error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	if exorderEntity == nil {
		// parse tx index from result
		n := parseIndexFromResult(result)
		exorderEntity = new(model.ExOrder)
		exorderEntity.From = result.From
		exorderEntity.To = result.To
		exorderEntity.Hash = result.Hash
		if n > 0 {
			exorderEntity.UUHash = fmt.Sprintf("%s:%s:%d", result.CoinType, result.Hash, n)
		} else {
			exorderEntity.UUHash = fmt.Sprintf("%s:%s", result.CoinType, result.Hash)
		}

		exorderEntity.Index = n
		exorderEntity.JadepoolOrderID = uint(jpOrderID)
		exorderEntity.Status = result.State
		exorderEntity.Type = result.BizType
		//todo:
		exorderEntity.AssetID = 1
		//todo:
		exorderEntity.JadepoolID = 1
		amount, condition, err := apd.NewFromString(result.Value)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.NewFromString error: %v, condition: %s", err, condition.String())
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
		exorderEntity.Amount = amount
		err = exorderRepo.Create(exorderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create exorder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	} else {
		if exorderEntity.Status == result.State {
			// repeat request
			utils.Infof("repeat request, jadepool order id: %d", exorderEntity.JadepoolOrderID)
			tx.Commit()
			resp := utils.Message(true, "success")
			utils.Respond(w, resp)
			return
		}

		updateEntity := &model.ExOrder{}
		updateEntity.Status = result.State
		exorderEntity.Status = result.State
		err = exorderRepo.UpdateColumns(exorderEntity.ID, updateEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("Update exorder error: %v", err)
			utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
			return
		}
	}

	// exorder has done, so create order
	if exorderEntity.Status == model.ExorderStatusDone {
		orderRepo := order.NewRepo(tx)
		orderEntity := new(model.Order)
		orderEntity.From = exorderEntity.From
		orderEntity.Hash = exorderEntity.Hash
		orderEntity.Index = exorderEntity.Index
		orderEntity.Status = model.OrderStatusInit
		orderEntity.To = exorderEntity.To
		orderEntity.Type = exorderEntity.Type
		orderEntity.UUHash = exorderEntity.UUHash
		orderEntity.AssetID = exorderEntity.AssetID
		orderEntity.Amount = exorderEntity.Amount
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
