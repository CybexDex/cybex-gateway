package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

const (
	// CybOrderStatusInit ...
	CybOrderStatusInit = "INIT"
	// CybOrderStatusHolding ...
	CybOrderStatusHolding = "HOLDING"
	// CybOrderStatusPending ...
	CybOrderStatusPending = "PENDING"
	// CybOrderStatusDone ...
	CybOrderStatusDone = "DONE"
	// CybOrderStatusFailed ...
	CybOrderStatusFailed = "FAILED"
)

//CybOrder ...
type CybOrder struct {
	gorm.Model

	AssetID uint `gorm:"not null" json:"assetID"` // 1 to n
	AppID   uint `gorm:"not null" json:"appID"`   // 1 to n

	// Accounting      Accounting `gorm:"foreignkey:AccountingRefer" json:"accounting"` // 1 to 1
	// AccountingRefer uint       `json:"accountingRefer"`

	From   string       `gorm:"type:varchar(128)" json:"from"`
	To     string       `gorm:"type:varchar(128)" json:"to"`
	Amount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	CybFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"cybFee"` // cybex blockchain fee
	// FeeAssetID uint

	Hash   string `gorm:"index;type:varchar(128);not null" json:"hash"`
	UUHash string `gorm:"type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status string `gorm:"type:varchar(32);not null" json:"status"`  // INIT, HOLDING, PENDING, DONE, FAILED
	Type   string `gorm:"type:varchar(32);not null" json:"type"`    // DEPOSIT, WITHDRAW, RECHARGE, SWEEP, FEESETTLE
}

//UpdateColumns ...
func (a *CybOrder) UpdateColumns(b *CybOrder) error {
	return GetDB().Model(CybOrder{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *CybOrder) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *CybOrder) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *CybOrder) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
