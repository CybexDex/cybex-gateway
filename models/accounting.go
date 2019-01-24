package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Accounting ...
type Accounting struct {
	gorm.Model

	// AssetID       uint         `gorm:"assetID"` // asset pay for blockchain fee
	BlockchainFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"blockchainFee"`
	InAmount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"inAmount"`
	OutAmount     *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"outAmount"`
	Fee           *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`
	Settled       bool         `gorm:"default:false" json:"settled"`
}

//UpdateColumns ...
func (a *Accounting) UpdateColumns(b *Accounting) error {
	return GetDB().Model(Accounting{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Accounting) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Accounting) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Accounting) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
