package types

import ecc "cybex-gateway/utils/ecc"

// JPOrderResult ...
type JPOrderResult struct {
	ID            string      `json:"id"`
	CoinName      string      `json:"coinName"`
	Txid          string      `json:"txid"`
	Meta          interface{} `json:"meta"`
	State         string      `json:"state"`
	BizType       string      `json:"bizType"`
	Type          string      `json:"type"`
	SubType       string      `json:"subType"`
	CoinType      string      `json:"coinType"`
	To            string      `json:"to"`
	Value         string      `json:"value"`
	Sequence      uint        `json:"sequence"`
	Confirmations int         `json:"confirmations"`
	CreateAt      int64       `json:"create_at"`
	UpdateAt      int64       `json:"update_at"`
	From          string      `json:"from"`
	N             int         `json:"n"`
	Fee           string      `json:"fee"`
	Hash          string      `json:"hash"`
	Block         int         `json:"block"`
	ExtraData     string      `json:"extraData"`
	Memo          string      `json:"memo"`
	SendAgain     bool        `json:"sendAgain"`
}

// JPEvent ...
type JPEvent struct {
	Code      int                    `json:"code"`
	Status    int                    `json:"status"`
	Message   string                 `json:"message"`
	Crypto    string                 `json:"crypto"`
	Timestamp int64                  `json:"timestamp"`
	Sig       *ecc.ECCSig            `json:"sig"`
	Result    map[string]interface{} `json:"result"`
}

// JPAddressResult ...
type JPAddressResult struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

// JPWithdrawResult ...
type JPWithdrawResult JPOrderResult

// JPAddressRequest ...
type JPAddressRequest struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Callback  string `json:"callback,omitempty"`
}

// JPWithdrawRequest ...
type JPWithdrawRequest struct {
	// Type      string `json:"type"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Sequence  uint   `json:"sequence"`
	Timestamp int64  `json:"timestamp"`
}

// JPSendData ...
type JPSendData struct {
	Crypto    string      `json:"crypto"`
	Hash      string      `json:"hash"`
	Encode    string      `json:"encode"`
	AppID     string      `json:"appid"`
	Timestamp int64       `json:"timestamp"`
	Sig       *ecc.ECCSig `json:"sig"`
	Data      interface{} `json:"data"`
}
