package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

// Asset ...
type Asset struct {
	gorm.Model
	Name         string `gorm:"unique" json:"name"`
	Blockchain   string `json:"blockchain"`
	CYBName      string `json:"cybname"`
	Confirmation string `json:"confirmation"`

	SmartContract  string `json:"smartContract"`
	GatewayAccount string `json:"gatewayAccount"`
	WithdrawPrefix string `json:"withdrawPrefix"`

	DepositSwitch  bool `json:"depositSwitch"`
	WithdrawSwitch bool `json:"withdrawSwitch"`

	MinDeposit  decimal.Decimal `gorm:"type:numeric" json:"minDeposit"`
	MinWithdraw decimal.Decimal `gorm:"type:numeric" json:"minWithdraw"`
	WithdrawFee decimal.Decimal `gorm:"type:numeric" json:"withdrawFee"`
	DepositFee  decimal.Decimal `gorm:"type:numeric" json:"depositFee"`

	Precision string `json:"precision"`
	ImgURL    string `json:"imgURL"`
	HashLink  string `json:"hashLink"`
}

// AssetsAll ...
func AssetsAll() (out []*Asset, err error) {
	err = db.Find(&out).Error
	return out, err
}