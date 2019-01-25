package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

const (
	// JPOrderStatusPending ...
	JPOrderStatusPending = "PENDING"
	// JPOrderStatusDone ...
	JPOrderStatusDone = "DONE"
	// JPOrderStatusFailed ...
	JPOrderStatusFailed = "FAILED"
	// JPOrderTypeDeposit ...
	JPOrderTypeDeposit = "DEPOSIT"
	// JPOrderTypeWithdraw ...
	JPOrderTypeWithdraw = "WITHDRAW"
)

//JPOrder ...
type JPOrder struct {
	gorm.Model

	AssetID    uint `gorm:"not null" json:"assetID"`    // n to 1
	JadepoolID uint `gorm:"not null" json:"jadepoolID"` // n to 1
	AppID      uint `gorm:"not null" json:"appID"`      // n to 1

	JadepoolOrderID uint         `gorm:"index;unique" json:"jadepoolOrderID"`
	From            string       `gorm:"type:varchar(128)" json:"from"`
	To              string       `gorm:"type:varchar(128)" json:"to"`
	Amount          *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Index           int          `json:"index"`
	Hash            string       `gorm:"index;type:varchar(128);not null" json:"hash"`
	UUHash          string       `gorm:"type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status          string       `gorm:"type:varchar(32);not null" json:"status"`  // PENDING, DONE, FAILED
	Type            string       `gorm:"type:varchar(32);not null" json:"type"`    // DEPOSIT, WITHDRAW
	Settled         bool         `gorm:"not null;default:false" json:"settled"`    // if order is created and count amount to balance, then Settled = true
}

//UpdateColumns ...
func (a *JPOrder) UpdateColumns(b *JPOrder) error {
	return GetDB().Model(JPOrder{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *JPOrder) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *JPOrder) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *JPOrder) Delete() (err error) {
	return GetDB().Delete(&a).Error
}

//AfterSave ... will be called each time after CREATE / SAVE / UPDATE
func (a JPOrder) AfterSave(tx *gorm.DB) (err error) {
	if a.Type == JPOrderTypeDeposit {
		if a.Settled == false {
			// set JPOrder Settled = true and SAVE to DB

			if a.Status == JPOrderStatusDone {
				// DEPOSIT JPOrder NOT settled before
				// status: -> DONE
				// balance: InLock -= 0, balance += amount
				// create ORDER

			} else if a.Status == JPOrderStatusFailed {
				// DEPOSIT JPOrder NOT settled before
				// status: -> FAILED
				// balance: InLock -= 0, balance -= 0
				// do NOTHING

			} else if a.Status == JPOrderStatusPending {
				// DEPOSIT JPOrder NOT settled before
				// status: -> PENDING
				// balance: InLock += amount, balance += 0

			}
		} else if a.Settled {
			if a.Status == JPOrderStatusDone {
				// DEPOSIT JPOrder settled before
				// status: PENDING -> DONE
				// balance: InLock -= 0, balance += 0
				// do NOTHING

			} else if a.Status == JPOrderStatusFailed {
				// DEPOSIT JPOrder settled before
				// status: PENDING -> FAILED
				// balance: InLock -= amount, balance += 0

			} else if a.Status == JPOrderStatusPending {
				// DEPOSIT JPOrder settled before
				// status is still PENDING
				// do NOTHING

			}
		}
	} else if a.Type == JPOrderTypeWithdraw {
		if a.Status == JPOrderStatusDone {
			// WITHDRAW JPOrder NOT settled before
			// status: -> DONE
			// balance: OutLock -= amount, balance -= 0

		} else if a.Status == JPOrderStatusFailed {
			// WITHDRAW JPOrder NOT settled before
			// status: -> FAILED
			// balance: OutLock -= amount, balance += amount ??
			// create NEW jporder, set it to order, move old jporder to order's

		} else if a.Status == JPOrderStatusPending {
			// status: -> PENDING
			// do NOTHING
		}
	}

	// if a.Status != JPOrderStatusDone {
	// 	u.Debugln("from JPOrder after save hook and the order status is DONE")
	// }

	// return errors.New("test error for rollback")

	return nil
}
