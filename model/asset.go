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
	JadeName     string `json:"-"`
	Blockchain   string `json:"blockchain"`
	Projectname  string `json:"projectname"`
	CYBName      string `json:"cybname"`
	CYBID        string `json:"cybid"`
	Confirmation string `json:"confirmation"`

	SmartContract  string `json:"smartContract"`
	GatewayAccount string `json:"gatewayAccount"` // 以account为准
	GatewayID      string `json:"-"`              // 方便读块时排除不需要的部分
	GatewayPass    string `json:"-"`              // 非常重要
	WithdrawPrefix string `json:"withdrawPrefix"`

	DepositSwitch  bool `json:"depositSwitch"`
	WithdrawSwitch bool `json:"withdrawSwitch"`

	MinDeposit  decimal.Decimal `gorm:"type:numeric" json:"minDeposit"`
	MinWithdraw decimal.Decimal `gorm:"type:numeric" json:"minWithdraw"`
	WithdrawFee decimal.Decimal `gorm:"type:numeric" json:"withdrawFee"`
	DepositFee  decimal.Decimal `gorm:"type:numeric" json:"depositFee"`

	CybAutoValue decimal.Decimal `gorm:"type:numeric" json:"cybAutoValue"`

	Precision string         `json:"precision"`
	ImgURL    string         `json:"imgURL"`
	HashLink  string         `json:"hashLink"`
	Info      postgres.Jsonb `gorm:"default:'{}'" json:"info"`
	UseMemo   bool           `json:"useMemo"`
	Disabled  bool           `json:"-"`
}

// Save ...
func (j *Asset) Save() error {
	return db.Save(j).Error
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
	err = db.Where(&Asset{
		Name: asset,
	}).Where(disableStr).First(&out).Error
	return out, err
}

// AssetsFrist ...
func AssetsFrist(query *Asset) (out *Asset, err error) {
	out = &Asset{}
	err = db.Where(query).Where(disableStr).First(&out).Error
	return out, err
}

// AssetsAllTest ...
func AssetsAllTest() (out []*Asset, err error) {
	err = db.Find(&out).Error
	return out, err
}

// AssetsFindTest ...
func AssetsFindTest(asset string, isPro bool) (out *Asset, err error) {
	out = &Asset{}
	err = db.Where(&Asset{
		Name: asset,
	}).First(&out).Error
	return out, err
}

// AssetsFristTest ...
func AssetsFristTest(query *Asset, isPro bool) (out *Asset, err error) {
	out = &Asset{}
	err = db.Where(query).First(&out).Error
	return out, err
}
