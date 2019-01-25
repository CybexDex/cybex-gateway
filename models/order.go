package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

const (
	// OrderStatusInit ...
	OrderStatusInit = "INIT"
	// OrderStatusProcessing ...
	OrderStatusProcessing = "PROCESSING"
	// OrderStatusDone ...
	OrderStatusDone = "DONE"
	// OrderStatusTerminated ...
	OrderStatusTerminated = "TERMINATED"
	// OrderTypeDeposit ...
	OrderTypeDeposit = "DEPOSIT"
	// OrderTypeWithdraw ...
	OrderTypeWithdraw = "WITHDRAW"
)

//Order ...
type Order struct {
	gorm.Model

	JPOrderID  uint     `json:"jpOrderID"`
	JPOrder    JPOrder  `json:"jpOrder"`
	CybOrderID uint     `json:"cybOrderID"`
	CybOrder   CybOrder `json:"cybOrder"`

	AssetID uint `gorm:"not null" json:"assetID"` // 1 to n
	AppID   uint `gorm:"not null" json:"appID"`   // 1 to n

	// Accounting      Accounting `gorm:"foreignkey:AccountingRefer" json:"accounting"` // 1 to 1
	// AccountingRefer uint       `json:"accountingRefer"`

	// Amount        *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	// BlockchainFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"blockchainFee"` //
	// Settled     bool         `gorm:"default:false" json:"settled"`

	JPHash   string `gorm:"index;type:varchar(128)" json:"jpHash"`
	JPUUHash string `gorm:"type:varchar(256)" json:"jpUUHash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	CYBHash   string `gorm:"index;type:varchar(128)" json:"cybHash"`
	CYBUUHash string `gorm:"type:varchar(256)" json:"cybUUHash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10)" json:"amount"`
	Fee         *apd.Decimal `gorm:"type:numeric(30,10)" json:"fee"` // fee in Asset

	Status string `gorm:"type:varchar(32);not null" json:"status"` // INIT, PROCESSING, DONE, TERMINATED
	Type   string `gorm:"type:varchar(32);not null" json:"type"`   // DEPOSIT, WITHDRAW
}

//UpdateColumns ...
func (a *Order) UpdateColumns(b *Order) error {
	return GetDB().Model(Order{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Order) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Order) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Order) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
