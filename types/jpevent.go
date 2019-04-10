package types

// JPOrderResult ...
type JPOrderResult struct {
	ID            string      `json:"id"`
	CoinName      string      `json:"coinName"`
	Txid          string      `json:"txid"`
	Meta          interface{} `json:"meta"`
	State         string      `json:"state"`
	BizType       string      `json:"bizType"`
	Type          string      `json:"type"`
	CoinType      string      `json:"coinType"`
	To            string      `json:"to"`
	Value         string      `json:"value"`
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
	Code      int    `json:"code"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Crypto    string `json:"crypto"`
	Timestamp int64  `json:"timestamp"`
	Sig       struct {
		R string `json:"r"`
		S string `json:"s"`
		V int    `json:"v"`
	} `json:"sig"`
	Result JPOrderResult `json:"result"`
}
