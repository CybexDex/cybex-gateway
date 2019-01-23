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

	"github.com/jinzhu/gorm"

	"github.com/cockroachdb/apd"

	model "git.coding.net/bobxuyang/cy-gateway-BN/models"
	exevent "git.coding.net/bobxuyang/cy-gateway-BN/repository/exevent"
	exorder "git.coding.net/bobxuyang/cy-gateway-BN/repository/exorder"
	order "git.coding.net/bobxuyang/cy-gateway-BN/repository/order"
	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"
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
		w.WriteHeader(400)
		return
	}
	utils.Infof("order noti request:\n %s", requestBody)

	request := OrderNotiRequest{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		utils.Errorf("Unmarshal error: %v", err)
		w.WriteHeader(400)
		return
	}

	result := request.Result
	result.Timestamp = request.Timestamp

	// 签名验证
	pubKey := os.Getenv("pub_key")
	ok, err := verifySign(result, request.Sig, pubKey)
	if err != nil {
		utils.Errorf("verifySign error: %v", err)
		w.WriteHeader(400)
		return
	}
	if !ok {
		utils.Errorf("verify result: %v", ok)
		w.WriteHeader(400)
		return
	}

	// 开启事务处理
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
		w.WriteHeader(400)
		return
	}

	// find exorder and save/update exorder
	jpOrderID, err := strconv.Atoi(result.ID)
	if err != nil {
		tx.Rollback()
		utils.Errorf("atoi error: %v, id: %s", err, result.ID)
		w.WriteHeader(400)
		return
	}
	exorderRepo := exorder.NewRepo(tx)
	exorderEntity, err := exorderRepo.GetByJPID(uint(jpOrderID))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		utils.Errorf("get exorder error: %v", err)
		w.WriteHeader(400)
		return
	}

	if exorderEntity == nil {
		// btc类的，如果订单发送给多个地址，n放在data结构里，需要解析取出来
		n := parseIndexFromResult(result)
		exorderEntity = new(model.ExOrder)
		exorderEntity.From = result.From
		exorderEntity.To = result.To
		exorderEntity.Hash = result.Hash
		exorderEntity.UUHash = fmt.Sprintf("%s:%s:%d", result.CoinType, result.Hash, n)
		exorderEntity.Index = n
		exorderEntity.JadepoolOrderID = uint(jpOrderID)
		exorderEntity.Status = result.State
		exorderEntity.Type = result.BizType
		//todo: 查询实际的id
		exorderEntity.AssetID = 1
		//todo: 查询实际的id
		exorderEntity.JadepoolID = 1
		amount, condition, err := apd.NewFromString(result.Value)
		if err != nil || condition.Any() {
			tx.Rollback()
			utils.Errorf("apd.NewFromString error: %v, condition: %s", err, condition.String())
			w.WriteHeader(400)
			return
		}
		exorderEntity.Amount = amount
		err = exorderRepo.Create(exorderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create exorder error: %v", err)
			w.WriteHeader(400)
			return
		}
	} else {
		if exorderEntity.Status == result.State {
			// 重复请求，返回正常结果，不执行后面创建order的操作
			utils.Infof("repeat request, jadepool order id: %d", exorderEntity.JadepoolOrderID)
			tx.Commit()
			resp := utils.Message(true, "success")
			utils.Respond(w, resp)
			return
		}

		// todo: save state
		exorderEntity.Status = result.State
	}

	// 链上已经确认，可以创建order
	if exorderEntity.Status == "done" {
		orderRepo := order.NewRepo(tx)
		orderEntity := new(model.Order)
		orderEntity.From = exorderEntity.From
		orderEntity.Hash = exorderEntity.Hash
		orderEntity.Index = exorderEntity.Index
		orderEntity.Status = "done"
		orderEntity.To = exorderEntity.To
		orderEntity.Type = exorderEntity.Type
		orderEntity.UUHash = exorderEntity.UUHash
		orderEntity.AssetID = exorderEntity.AssetID
		orderEntity.Amount = exorderEntity.Amount
		//todo: 查询实际的id
		orderEntity.AppID = 1
		err = orderRepo.Create(orderEntity)
		if err != nil {
			tx.Rollback()
			utils.Errorf("create order error: %v", err)
			w.WriteHeader(400)
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
	if result.CoinType == "BTC" {
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
	}

	return n
}
