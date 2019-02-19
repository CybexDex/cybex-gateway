package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/spf13/viper"
)

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
	/*strAppID := query.Get("appID")
	appID, err := strconv.Atoi(strAppID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	appRepo := app.NewRepo(model.GetDB())
	app, err := appRepo.GetByID(uint(appID))
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(app.JadepoolID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	utils.Debugln("jadepoolID", jadepool.ID)*/

	timestamp := time.Now().Unix() * 1000
	requestAddress := JPAddressRequest{}
	requestAddress.Timestamp = timestamp
	requestAddress.Type = coinType
	requestAddress.Callback = viper.GetString("jpsrv.self_addr") + "/api/order/noti"

	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = viper.GetString("jpsrv.jadepool_appid")
	sendData.Data = &requestAddress

	prikey := viper.GetString("jpsrv.pri_key")
	sig, err := utils.SignECCData(prikey, sendData.Data)
	if err != nil {
		utils.Errorf("SignECCData error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	sendData.Sig = sig

	bs, _ := json.Marshal(sendData)
	jadepoolAddr := viper.GetString("defaultJadepool.jadepool_addr")
	//jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	url := jadepoolAddr + "/api/v1/addresses/new"

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
		pubKey := viper.GetString("jadepool.pub_key")
		//pubKey := jadepool.EccPubKey
		data.Result["timestamp"] = data.Timestamp
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

// GetNewAddress2 ...
func GetNewAddress2(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	coinType := query.Get("type")
	/*strAppID := query.Get("appID")
	appID, err := strconv.Atoi(strAppID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	appRepo := app.NewRepo(model.GetDB())
	app, err := appRepo.GetByID(uint(appID))
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(app.JadepoolID)
	if err != nil {
		utils.Respond(w, utils.Message(false, "bad request"), http.StatusBadRequest)
		return
	}
	utils.Debugln("jadepoolID", jadepool.ID)*/

	callbackAddr := viper.GetString("jpsrv.self_addr")
	jadepoolAddr := viper.GetString("defaultJadepool.jadepool_addr")
	//jadepoolAddr := fmt.Sprintf("http://%s:%d", jadepool.Host, jadepool.Port)
	pubKey := viper.GetString("defaultJadepool.pub_key")
	//pubKey := jadepool.EccPubKey
	priKey := viper.GetString("jpsrv.pri_key")
	appID := "app"
	result, err := GetAddressFromJadepool(coinType, callbackAddr, appID, jadepoolAddr, priKey, pubKey)
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
