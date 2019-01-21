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
