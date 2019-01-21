package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Order ...
type Order struct {
	gorm.Model

	AssetID uint `gorm:"not null" json:"assetID"` // 1 to n
	AppID   uint `gorm:"not null" json:"appID"`   // 1 to n

	Accounting      Accounting `gorm:"foreignkey:AccountingRefer" json:"accounting"` // 1 to 1
	AccountingRefer uint       `json:"accountingRefer"`

	From   string       `gorm:"type:varchar(128)" json:"from"`
	To     string       `gorm:"type:varchar(128)" json:"to"`
	Amount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Index  int          `json:"index"`
	Hash   string       `gorm:"index;type:varchar(128);not null" json:"hash"`
	UUHash string       `gorm:"type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status string       `gorm:"type:varchar(32);not null" json:"status"`  // INIT, PROCESSING, PROCESSED, HOLDING, PENDING, DONE, FAILED, CONTINUE, RETRY, TERMINATED
	Type   string       `gorm:"type:varchar(32);not null" json:"type"`    // DEPOSIT, WITHDRAW, RECHARGE, INTERNAL
}
