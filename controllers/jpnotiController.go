package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"reflect"

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
	// btc类的，如果订单发送给多个地址，n放在data结构里，需要解析取出来
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
	n = n + 0
	//todo: save log

	//todo: find exorder and save/update exorder

	//todo: response
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
