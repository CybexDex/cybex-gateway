package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Balance ...
type Balance struct {
	gorm.Model

	AppID   uint `json:"appID"`
	AssetID uint `json:"assetID"`

	Balance   *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"balance"`
	InLocked  *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"inLocked"`  // after order DONE, add to balance
	OutLocked *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"outLocked"` // deduct from balance and create WITHDRAWED ORDER
}
