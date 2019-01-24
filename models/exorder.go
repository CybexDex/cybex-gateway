package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

const (
	// ExorderStatusPending ...
	ExorderStatusPending = "PENDING"
	// ExorderStatusDone ...
	ExorderStatusDone = "DONE"
	// ExorderStatusFailed ...
	ExorderStatusFailed = "FAILED"
	// ExorderTypeDeposit ...
	ExorderTypeDeposit = "DEPOSIT"
	// ExorderTypeWithdraw ...
	ExorderTypeWithdraw = "WITHDRAW"
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
	Settled         bool         `gorm:"not null;default:false" json:"settled"`    // if order is created and count amount to balance, then Settled = true
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

//AfterSave ... will be called each time after CREATE / SAVE / UPDATE
func (a ExOrder) AfterSave(tx *gorm.DB) (err error) {
	if a.Settled == true {
		if a.Status == ExorderStatusDone {
			if a.Type == ExorderTypeDeposit {
				// DEPOSIT exorder settled before
				// status: PENDING -> DONE
				// balance: InLock -= amount, balance += amount

			} else if a.Type == ExorderTypeWithdraw {
				// WITHDRAW exorder settled before
				// status: PENDING -> DONE
				// balance: OutLock -= amount, balance -= 0

			}
		} else if a.Status == ExorderStatusFailed {
			if a.Type == ExorderTypeDeposit {
				// DEPOSIT exorder settled before
				// status: PENDING -> FAILED
				// balance: InLock -= amount, balance += 0

			} else if a.Type == ExorderTypeWithdraw {
				// WITHDRAW exorder settled before
				// status: PENDING -> FAILED
				// balance: OutLock -= amount, balance += amount

			}
		} else if a.Status == ExorderStatusPending {
			// exorder was settled before
			// status is still PENDING
			// then do nothing
		}
	} else if a.Settled == false {
		// set exorder Settled = true and SAVE to DB

		if a.Status == ExorderStatusDone {
			if a.Type == ExorderTypeDeposit {
				// DEPOSIT exorder NOT settled before
				// status: -> DONE
				// balance: InLock -= 0, balance += amount

			} else if a.Type == ExorderTypeWithdraw {
				// WITHDRAW exorder NOT settled before
				// status: -> DONE
				// balance: OutLock -= 0, balance -= amount

			}
		} else if a.Status == ExorderStatusFailed {
			if a.Type == ExorderTypeDeposit {
				// DEPOSIT exorder NOT settled before
				// status: -> FAILED
				// balance: InLock -= 0, balance += 0
				// do NOTHING
			} else if a.Type == ExorderTypeWithdraw {
				// WITHDRAW exorder NOT settled before
				// status: -> FAILED
				// balance: OutLock -= 0, balance += 0
				// do NOTHING
			}
		} else if a.Status == ExorderStatusPending {
			if a.Type == ExorderTypeDeposit {
				// DEPOSIT exorder NOT settled before
				// status: -> PENDING
				// balance: InLock += amount, balance += 0

			} else if a.Type == ExorderTypeWithdraw {
				// WITHDRAW exorder NOT settled before
				// status: -> PENDING
				// balance: OutLock += amount, balance -= amount

			}
		}
	}

	// if a.Status != ExorderStatusDone {
	// 	u.Debugln("from exorder after save hook and the order status is DONE")
	// }

	return nil

	// return errors.New("test error for rollback")
}
