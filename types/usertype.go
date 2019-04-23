package types

import (
	"time"
)

// UserResultAddress ...
type UserResultAddress struct {
	Address  string    `json:"address"`
	Asset    string    `json:"asset"`
	CreateAt time.Time `json:"createAt"`
	CybName  string    `json:"cybName"`
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
