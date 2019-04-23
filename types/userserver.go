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
