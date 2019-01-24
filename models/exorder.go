package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"

	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

const (
	// ExorderStatusPending ...
	ExorderStatusPending = "PENDING"
	// ExorderStatusDone ...
	ExorderStatusDone = "DONE"
	// ExorderStatusFailed ...
	ExorderStatusFailed = "FAILED"
)

//ExOrder ...
type ExOrder struct {
	gorm.Model

	AssetID    uint `gorm:"not null" json:"assetID"`    // n to 1
	JadepoolID uint `gorm:"not null" json:"jadepoolID"` // n to 1

	JadepoolOrderID uint         `gorm:"index;unique" json:"jadepoolOrderID"`
	From            string       `gorm:"type:varchar(128)" json:"from"`
	To              string       `gorm:"type:varchar(128)" json:"to"`
	Amount          *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Index           int          `json:"index"`
	Hash            string       `gorm:"index;type:varchar(128);not null" json:"hash"`
	UUHash          string       `gorm:"type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status          string       `gorm:"type:varchar(32);not null" json:"status"`  // PENDING, DONE, FAILED
	Type            string       `gorm:"type:varchar(32);not null" json:"type"`    // DEPOSIT, WITHDRAW
}

//UpdateColumns ...
func (a *ExOrder) UpdateColumns(b *ExOrder) error {
	return GetDB().Model(ExOrder{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *ExOrder) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *ExOrder) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *ExOrder) Delete() (err error) {
	return GetDB().Delete(&a).Error
}

//AfterSave ...
func (a ExOrder) AfterSave(tx *gorm.DB) (err error) {
	if a.Status == "DONE" {
		u.Debugln("from exorder after save hook and the order status is DONE")
	}

	return nil

	// return errors.New("test error for rollback")
}
