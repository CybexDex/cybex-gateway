package jpsrv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jadepool"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/spf13/viper"
)

// JPAddressRequest ...
type JPAddressRequest struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Callback  string `json:"callback,omitempty"`
}

// RequestNewAddress ...
func RequestNewAddress(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	coinType := query.Get("type")
	strJadepoolID := query.Get("jadepoolID")
	if len(strJadepoolID) == 0 {
		strJadepoolID = "1"
	}
	jadepoolID, err := strconv.Atoi(strJadepoolID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(uint(jadepoolID))
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}

	callbackAddr := viper.GetString("jpsrv.self_addr")
	jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	pubKey := jadepool.EccPubKey
	priKey := viper.GetString("jpsrv.pri_key")
	jadepoolAppID := viper.GetString("jpsrv.jadepool_appid")
	result, err := GetAddressFromJadepool(coinType, callbackAddr, jadepoolAppID, jadepoolAddr, priKey, pubKey)
	if err != nil {
		utils.Errorf("err: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", result)
	utils.Respond(w, resp)
}

// VerifyBlockchainAddress ...
func VerifyBlockchainAddress(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	coinType := query.Get("type")
	strJadepoolID := query.Get("jadepoolID")
	addr := query.Get("addr")
	if len(strJadepoolID) == 0 {
		strJadepoolID = "1"
	}
	jadepoolID, err := strconv.Atoi(strJadepoolID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(uint(jadepoolID))
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}

	jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	pubKey := jadepool.EccPubKey
	priKey := viper.GetString("jpsrv.pri_key")
	jadepoolAppID := viper.GetString("jpsrv.jadepool_appid")
	result, err := VerifyAddress(coinType, addr, jadepoolAppID, jadepoolAddr, priKey, pubKey)
	if err != nil {
		utils.Errorf("err: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", result)
	utils.Respond(w, resp)
}

// RequestConfirmations ...
func RequestConfirmations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	coinType := query.Get("type")
	strJadepoolID := query.Get("jadepoolID")
	if len(strJadepoolID) == 0 {
		strJadepoolID = "1"
	}
	jadepoolID, err := strconv.Atoi(strJadepoolID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(uint(jadepoolID))
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}

	jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	pubKey := jadepool.EccPubKey
	priKey := viper.GetString("jpsrv.pri_key")
	jadepoolAppID := viper.GetString("jpsrv.jadepool_appid")
	result, err := GetCoinConfirmations(coinType, jadepoolAppID, jadepoolAddr, priKey, pubKey)
	if err != nil {
		utils.Errorf("err: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", result)
	utils.Respond(w, resp)
}

// GetAddressFromJadepool ...
func GetAddressFromJadepool(coinType, callbackAddr, appID, jadepoolAddr string, priKey string, pubKey string) (map[string]interface{}, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coin type is empty")
	}
	coinType = strings.ToUpper(coinType)

	timestamp := time.Now().Unix() * 1000
	requestAddress := JPAddressRequest{}
	requestAddress.Timestamp = timestamp
	requestAddress.Type = coinType
	requestAddress.Callback = callbackAddr + "/api/order/noti"

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = appID
	sendData.Data = &requestAddress

	sig, err := utils.SignECCData(priKey, sendData.Data)
	if err != nil {
		return nil, fmt.Errorf("SignECCData error: %v", err)
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	url := jadepoolAddr + "/api/v1/addresses/new"

	data := JPComeData{}
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	if os.Getenv("env") != "dev" {
		// verify sig
		data.Result["timestamp"] = data.Timestamp
		ok, err := utils.VerifyECCSign(data.Result, &data.Sig, pubKey)
		if err != nil {
			return nil, fmt.Errorf("verifySign error: %v, data: %#v", err, data)
		}
		if !ok {
			return nil, fmt.Errorf("verify result: %v, data: %#v", ok, data)
		}
	}

	if data.Code != 0 || data.Status != 0 || data.Result == nil || data.Result["address"] == nil {
		return nil, fmt.Errorf("not found address, data: %#v", data)
	}

	return data.Result, nil
}

// VerifyAddress ...
func VerifyAddress(coinType, addr, appID, jadepoolAddr string, priKey string, pubKey string) (map[string]interface{}, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coin type is empty")
	}
	coinType = strings.ToUpper(coinType)

	timestamp := time.Now().Unix() * 1000
	requestAddress := JPAddressRequest{}
	requestAddress.Type = coinType

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = appID
	sendData.Data = &requestAddress

	sig, err := utils.SignECCData(priKey, sendData.Data)
	if err != nil {
		return nil, fmt.Errorf("SignECCData error: %v", err)
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	url := jadepoolAddr + "/api/v1/addresses/" + addr + "/verify"

	data := JPComeData{}
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	if os.Getenv("env") != "dev" {
		// verify sig
		data.Result["timestamp"] = data.Timestamp
		ok, err := utils.VerifyECCSign(data.Result, &data.Sig, pubKey)
		if err != nil {
			return nil, fmt.Errorf("verifySign error: %v, data: %#v", err, data)
		}
		if !ok {
			return nil, fmt.Errorf("verify result: %v, data: %#v", ok, data)
		}
	}

	if data.Code != 0 || data.Status != 0 || data.Result == nil {
		return nil, fmt.Errorf("not found result, data: %#v", data)
	}

	return data.Result, nil
}

// GetCoinConfirmations ...
func GetCoinConfirmations(coinType, appID, jadepoolAddr string, priKey string, pubKey string) (map[string]interface{}, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coin type is empty")
	}
	coinType = strings.ToUpper(coinType)

	result := make(map[string]interface{})
	result["type"] = coinType
	result["confirmations"] = 30
	return result, nil

	/*timestamp := time.Now().Unix() * 1000
	requestAddress := JPAddressRequest{}
	requestAddress.Type = coinType

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = appID
	sendData.Data = &requestAddress

	sig, err := utils.SignECCData(priKey, sendData.Data)
	if err != nil {
		return nil, fmt.Errorf("SignECCData error: %v", err)
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	url := jadepoolAddr + "/api/v1/confirmations"

	data := JPComeData{}
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	if os.Getenv("env") != "dev" {
		// verify sig
		data.Result["timestamp"] = data.Timestamp
		ok, err := utils.VerifyECCSign(data.Result, &data.Sig, pubKey)
		if err != nil {
			return nil, fmt.Errorf("verifySign error: %v, data: %#v", err, data)
		}
		if !ok {
			return nil, fmt.Errorf("verify result: %v, data: %#v", ok, data)
		}
	}

	if data.Code != 0 || data.Status != 0 || data.Result == nil {
		return nil, fmt.Errorf("not found result, data: %#v", data)
	}

	return data.Result, nil*/
}
