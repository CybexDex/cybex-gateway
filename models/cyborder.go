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
	// CybOrderTypeDeposit ...
	CybOrderTypeDeposit = "DEPOSIT"
	// CybOrderTypeWithdraw ...
	CybOrderTypeWithdraw = "WITHDRAW"
	// CybOrderTypeRecharge ...
	CybOrderTypeRecharge = "RECHARGE"
	// CybOrderTypeSweep ...
	CybOrderTypeSweep = "SWEEP"
	// CybOrderTypeFeeSettle ...
	CybOrderTypeFeeSettle = "FEESETTLE"
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
	CybFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"cybFee"` // cybex blockchain fee
	// FeeAssetID uint

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"` // fee in Asset

	Hash    string `gorm:"unique;index;type:varchar(128);not null" json:"hash"`
	UUHash  string `gorm:"unique;nidex;type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status  string `gorm:"type:varchar(32);not null" json:"status"`               // INIT, HOLDING, PENDING, DONE, FAILED
	Type    string `gorm:"type:varchar(32);not null" json:"type"`                 // DEPOSIT, WITHDRAW, RECHARGE, SWEEP, FEESETTLE
	Settled bool   `gorm:"not null;default:false" json:"settled"`                 // if order is created and count amount to balance, then Settled = true
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

//AfterSave ... will be called each time after CREATE / SAVE / UPDATE
func (a CybOrder) AfterSave(tx *gorm.DB) (err error) {
	if a.Type == CybOrderTypeWithdraw {
		if a.Settled == false {
			// case 1, 2, 4, 6 will executed only once

			// set CybOrder Settled = true and SAVE to DB

			if a.Status == CybOrderStatusDone { // case 1
				// DEPOSIT CybOrder NOT settled before
				// status: -> DONE
				// balance: OutLock += cyborder.amount, balance -= cyborder.amount, same as case 3
				// create ORDER, order.TotalAmount = cyborder.amount

			} else if a.Status == CybOrderStatusFailed { // case 2
				// DEPOSIT CybOrder NOT settled before
				// status: -> FAILED
				// balance: OutLock -= 0, balance -= 0
				// do NOTHING

			} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding { // case 3
				// DEPOSIT CybOrder NOT settled before
				// status: -> PENDING
				// balance: OutLock += cyborder.amount, balance -= cyborder.amount, same as case 1

			}
		} else if a.Settled {
			if a.Status == CybOrderStatusDone { // case 4
				// DEPOSIT CybOrder settled before
				// status: PENDING -> DONE
				// balance: OutLock -= 0, balance += 0
				// create ORDER, order.TotalAmount = cyborder.amount

			} else if a.Status == CybOrderStatusFailed { // case 5, symmetrical to case 3 & 1
				// DEPOSIT CybOrder settled before
				// status: PENDING -> FAILED
				// balance: OutLock -= cyborder.amount, balance += cyborder.amount

			} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding { // case 6
				// DEPOSIT CybOrder settled before
				// status: PENDING -> PENDING
				// do NOTHING

			}
		}
	} else if a.Type == CybOrderTypeDeposit { // map to jporder's part
		if a.Status == CybOrderStatusDone {
			// WITHDRAW CybOrder NOT settled before
			// status: -> DONE
			// balance: balance = balance + order.Amount, inLocked -= order.TotalAmount, inLockedFee -= order.Fee

		} else if a.Status == CybOrderStatusFailed {
			// WITHDRAW CybOrder NOT settled before
			// status: -> FAILED
			// balance: InLock -= 0, balance -= 0, inLockedFee -=0
			// create NEW cyborder, set it to order, move old cyborder to order's FailedCybOrders

		} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding {
			// status: -> PENDING
			// do NOTHING
		}
	}

	// return errors.New("test error for rollback")

	return nil
}
