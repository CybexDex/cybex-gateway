package types

import (
	"coding.net/yundkyy/cybexgolib/operations"
	cybTypes "coding.net/yundkyy/cybexgolib/types"
)

// GatewayAccount ...
type GatewayAccount struct {
	Account *cybTypes.Account
	Type    string
	Asset   string
}

// HandleInterface ...
type HandleInterface interface {
	HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction)
}

// AssetConfig ...
type AssetConfig struct {
	Name         string `json:"name"`
	HandleAction string `json:"handle_action"`
	Deposit      struct {
		Gateway     string   `json:"gateway"`
		Gatewaypass string   `json:"gatewaypass"`
		Sendto      []string `json:"sendto"`
		Switch      bool     `json:"switch"`
		JustAsset   string   `json:"just_asset"`
		Just        []string `json:"just"`
	} `json:"deposit"`
	Withdraw struct {
		Gateway     string `json:"gateway"`
		Gatewaypass string `json:"gatewaypass"`
		Coin        string `json:"coin"`
		Wait        string `json:"wait"`
		Send        string `json:"send"`
		Switch      bool   `json:"switch"`
	} `json:"withdraw"`
}
