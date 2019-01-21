package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//ExOrder ...
type ExOrder struct {
	gorm.Model

	AssetID    uint `gorm:"not null" json:"assetID"`    // n to 1
	JadepoolID uint `gorm:"not null" json:"jadepoolID"` // n to 1

	From   string       `gorm:"type:varchar(128)" json:"from"`
	To     string       `gorm:"type:varchar(128)" json:"to"`
	Amount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Index  int          `json:"index"`
	Hash   string       `gorm:"index;type:varchar(128);not null" json:"hash"`
	UUHash string       `gorm:"type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status string       `gorm:"type:varchar(32);not null" json:"status"`  // PENDING, DONE, FAILED
	Type   string       `gorm:"type:varchar(32);not null" json:"type"`    // DEPOSIT, WITHDRAW
}
