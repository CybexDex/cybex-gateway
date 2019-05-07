package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserResultAddress ...
type UserResultAddress struct {
	Address  string    `json:"address"`
	Asset    string    `json:"asset"`
	CreateAt time.Time `json:"createAt"`
	CybName  string    `json:"cybName"`
}

// VerifyRes ...
type VerifyRes struct {
	Address   string `json:"address"`
	Asset     string `json:"asset"`
	Timestamp uint   `json:"timestamp"`
	Valid     bool   `json:"valid"`
}

// UserResultBBB ...
type UserResultBBB struct {
	Name            string `json:"name"`
	Blockchain      string `json:"blockchain"`
	DepositAs       string `json:"depositAs"`
	WithdrawAsset   string `json:"withdrawAsset"`
	WithdrawGateway string `json:"withdrawGateway"`
	WithdrawPrefix  string `json:"withdrawPrefix"`

	DepositSwitch  bool `json:"depositSwitch"`
	WithdrawSwitch bool `json:"withdrawSwitch"`

	MinDeposit  string `json:"minDeposit"`
	MinWithdraw string `json:"minWithdraw"`
	WithdrawFee string `json:"withdrawFee"`
	DepositFee  string `json:"depositFee"`
}

// UserResultAsset ...
type UserResultAsset struct {
	Name         string `json:"name"`
	Blockchain   string `json:"blockchain"`
	CYBName      string `json:"cybname"`
	Confirmation string `json:"confirmation"`

	SmartContract  string `json:"smartContract"`
	GatewayAccount string `json:"gatewayAccount"`
	WithdrawPrefix string `json:"withdrawPrefix"`

	DepositSwitch  bool `json:"depositSwitch"`
	WithdrawSwitch bool `json:"withdrawSwitch"`

	MinDeposit  decimal.Decimal `json:"minDeposit"`
	MinWithdraw decimal.Decimal `json:"minWithdraw"`
	WithdrawFee decimal.Decimal `json:"withdrawFee"`
	DepositFee  decimal.Decimal `json:"depositFee"`

	Precision uint   `json:"precision"`
	ImgURL    string `json:"imgURL"`
	HashLink  string `json:"hashLink"`
}

// RecordsQuery ...
type RecordsQuery struct {
	FundType string `form:"fundType"`
	LastID   string `form:"lastid"`
	Size     string `form:"size"`
	Asset    string `form:"asset"`
	User     string `form:"user"`
}

// Record ...
type Record struct {
	Type        string    `json:"type"`
	ID          uint      `json:"id"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CybexName   string    `json:"cybexName"`
	OutAddr     string    `json:"outAddr"`
	GatewayAddr string    `json:"gatewayAddr"`
	Confirms    string    `json:"confirms"`
	Asset       string    `json:"asset"`
	OutHash     string    `json:"outHash"`
	CybHash     *string   `json:"cybHash"`
	TotalAmount string    `json:"totalAmount"`
	Amount      string    `json:"amount"`
	Fee         string    `json:"fee"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	Link        string    `json:"link"`
}
