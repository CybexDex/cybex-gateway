package controllers

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestSignTransaction(t *testing.T) {
	prikey := "bf12996feeaa2977b6ca0d33a0e8bd2ccfc4844c6f8a7e6d15c099f8da4a255c"
	timestamp := time.Now().Unix() * 1000
	sendData := &JPSendData{}
	sendData.Crypto = "ecc"
	sendData.Encode = "base64"
	sendData.Timestamp = timestamp
	sendData.Hash = "sha3"
	sendData.AppID = "app"
	sendData.Data = &JPTransaction{}
	sendData.Data.Timestamp = timestamp

	sig, err := signTransaction(prikey, sendData.Data)
	if err != nil {
		t.Error(err)
		return
	}
	sendData.Sig = sig
	sendByts, err := json.Marshal(&sendData)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(sendByts))
}
