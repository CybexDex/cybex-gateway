package model

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shopspring/decimal"
)

// Asset ...
type Asset struct {
	gorm.Model
	Name         string `gorm:"unique" json:"name"`
	Blockchain   string `json:"blockchain"`
	Projectname  string `json:"projectname"`
	CYBName      string `json:"cybname"`
	CYBID        string `json:"cybid"`
	Confirmation string `json:"confirmation"`

	SmartContract  string `json:"smartContract"`
	GatewayAccount string `json:"gatewayAccount"`
	GatewayPass    string `json:"-"` // 非常重要
	WithdrawPrefix string `json:"withdrawPrefix"`

	DepositSwitch  bool `json:"depositSwitch"`
	WithdrawSwitch bool `json:"withdrawSwitch"`

	MinDeposit  decimal.Decimal `gorm:"type:numeric" json:"minDeposit"`
	MinWithdraw decimal.Decimal `gorm:"type:numeric" json:"minWithdraw"`
	WithdrawFee decimal.Decimal `gorm:"type:numeric" json:"withdrawFee"`
	DepositFee  decimal.Decimal `gorm:"type:numeric" json:"depositFee"`

	Precision string         `json:"precision"`
	ImgURL    string         `json:"imgURL"`
	HashLink  string         `json:"hashLink"`
	Info      postgres.Jsonb `gorm:"default:'{}'" json:"info"`
	UseMemo   bool           `json:"useMemo"`
	Disabled  bool           `json:"-"`
}

// TestJP ...
func (j *Asset) TestJP() error {
	return nil
}

// TestCybex ...
func (j *Asset) TestCybex() error {
	return nil
}

const disableStr = "disabled is NULL or disabled is false"

// AssetsAll ...
func AssetsAll() (out []*Asset, err error) {
	err = db.Find(&out, disableStr).Error
	return out, err
}

// AssetsFind ...
func AssetsFind(asset string) (out *Asset, err error) {
	out = &Asset{}
	err = db.First(&out, &Asset{
		Name: asset,
	}, disableStr).Error
	return out, err
}

// AssetsFrist ...
func AssetsFrist(query *Asset) (out *Asset, err error) {
	out = &Asset{}
	err = db.First(&out, query, disableStr).Error
	return out, err
}
