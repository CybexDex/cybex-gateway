package sass

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// API ...
type API struct {
	APIAddr    string `json:"apiAddr"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	nonceCount int
}

// APIResponse ...
type APIResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// AddressRequest ...
type AddressRequest struct {
	Address   string `json:"address,omitempty"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign,omitempty"`
}

// AddressResponse ...
type AddressResponse struct {
	Address int64  `json:"address"`
	Mode    string `json:"mode"`
}

// GetAdderess ...
func (a *API) GetAdderess(coinType string) (map[string]interface{}, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	url := fmt.Sprintf("%s/api/v1/address/%s/new", a.APIAddr, coinType)
	data := &AddressRequest{
		Timestamp: time.Now().Unix(),
		Nonce:     a.genNonce(),
	}
	sign, err := SignHMACSHA256(data, a.AppSecret)
	if err != nil {
		return nil, err
	}
	data.Sign = sign
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-App-Key", a.AppKey)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	respData := APIResponse{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil, err
	}

	if respData.Code != 0 {
		return nil, errors.New(respData.Message)
	}

	return respData.Data, nil
}

// VerifyAddress ...
func (a *API) VerifyAddress(coinType, addr string) (map[string]interface{}, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	url := fmt.Sprintf("%s/api/v1/address/%s/verify", a.APIAddr, coinType)
	data := &AddressRequest{
		Timestamp: time.Now().Unix(),
		Nonce:     a.genNonce(),
		Address:   addr,
	}
	sign, err := SignHMACSHA256(data, a.AppSecret)
	if err != nil {
		return nil, err
	}
	data.Sign = sign
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-App-Key", a.AppKey)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	respData := APIResponse{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil, err
	}

	if respData.Code != 0 {
		return nil, errors.New(respData.Message)
	}

	return respData.Data, nil
}

// WithdrawRequest ...
type WithdrawRequest struct {
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	ID        string `json:"id"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Sign      string `json:"sign,omitempty"`
}

// OrderInfo ...
type OrderInfo struct {
	ID                      string `json:"id"`
	From                    string `json:"from"`
	To                      string `json:"to"`
	Value                   string `json:"value"`
	Type                    string `json:"type"`
	CoinName                string `json:"coinName"`
	BizType                 string `json:"bizType"`
	Hash                    string `json:"txid"`
	Fee                     string `json:"fee"`
	Confirmations           int    `json:"confirmations"`
	N                       int    `json:"n"`
	Memo                    string `json:"memo"`
	State                   string `json:"state"`
	Block                   int    `json:"block"`
	Sign                    string `json:"sign,omitempty"`
	WithdrawID              string `json:"withdrawID"`
	AffirmativeConfirmation int    `json:"affirmativeConfirmation"`
}

// Withdraw ...
func (a *API) Withdraw(id, coinType, to, value, memo string) (*OrderInfo, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}
	url := fmt.Sprintf("%s/api/v1/app/%s/withdraw", a.APIAddr, coinType)
	data := &WithdrawRequest{
		Timestamp: time.Now().Unix(),
		Nonce:     a.genNonce(),
		ID:        id,
		To:        to,
		Value:     value,
	}
	sign, err := SignHMACSHA256(data, a.AppSecret)
	if err != nil {
		return nil, err
	}
	data.Sign = sign
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-App-Key", a.AppKey)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}

	respData := APIResponse{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil, err
	}
	if respData.Code != 0 {
		return nil, errors.New(respData.Message)
	}

	resultData, err := json.Marshal(respData.Data)
	if err != nil {
		return nil, err
	}
	result := OrderInfo{}
	err = json.Unmarshal(resultData, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ParseOrderNoti ...
func (a *API) ParseOrderNoti(data []byte) (*OrderInfo, error) {
	info := OrderInfo{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	sign := info.Sign
	if len(sign) == 0 {
		return nil, errors.New("sign empty")
	}
	info.Sign = ""

	mySign, err := SignHMACSHA256(info, a.AppSecret)
	if err != nil {
		return nil, err
	}

	if mySign != sign {
		return nil, errors.New("sign error")
	}

	return &info, nil
}

func (a *API) genNonce() string {
	a.nonceCount++
	timestamp := time.Now().UnixNano()
	rand := rand.Int63n(timestamp)
	return fmt.Sprintf("%d%d%d", a.nonceCount, timestamp, rand)
}
