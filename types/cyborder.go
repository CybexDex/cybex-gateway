package types

import (
	"github.com/CybexDex/cybex-go/operations"
	cybTypes "github.com/CybexDex/cybex-go/types"
)

// GatewayAccount ...
type GatewayAccount struct {
	Account *cybTypes.Account
	Type    string
	Asset   string
	MemoPri cybTypes.PrivateKeys
}

// HandleInterface ...
type HandleInterface interface {
	HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction, prefix string)
}

// AssetConfig ...
type AssetConfig struct {
	Name         string `json:"name"`
	BlockChain   string `json:"blockChain"`
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
		Fee         string `json:"fee"`
		Memopre     string `json:"memopre"`
	} `json:"withdraw"`
}
