package model

import (
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
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

	JadepoolOrderID uint   `gorm:"index;unique" json:"jadepoolOrderID"`
	From            string `gorm:"type:varchar(128)" json:"from"`
	To              string `gorm:"type:varchar(128)" json:"to"`

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"` // fee in Asset

	Index   int    `json:"index"`
	Hash    string `gorm:"unique;index;type:varchar(128);not null" json:"hash"`
	UUHash  string `gorm:"unique;index;type:varchar(256);not null" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	Status  string `gorm:"type:varchar(32);not null" json:"status"`               // INIT, PENDING, DONE, FAILED
	Type    string `gorm:"type:varchar(32);not null" json:"type"`                 // DEPOSIT, WITHDRAW
	Settled bool   `gorm:"not null;default:false" json:"settled"`                 // if count amount to balance, then Settled = true
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
func (a *JPOrder) AfterSave1(tx *gorm.DB) (err error) {
	if a.Type == JPOrderTypeDeposit {
		if a.Settled == false {
			// case 1, 2, 4, 6 will executed only once

			// set JPOrder Settled = true and SAVE to DB
			a.Settled = true
			err := tx.Save(a).Error
			if err != nil {
				u.Errorln("set order settled to true error", a.ID)
				return err
			}
			u.Debugln("set order settled to true", a.ID)

			if a.Status == JPOrderStatusDone { // case 1
				// DEPOSIT JPOrder NOT settled before
				// status: -> DONE
				// balance: InLock += jporder.amount, balance += 0, InLockedFee += fee from asset, same as case 3
				// create ORDER, order.TotalAmount = jporder.amount, order.Fee = fee from asset table, order.Amount = order.TotalAmount - order.Fee

				// update balance record
				var bal Balance
				err = tx.FirstOrCreate(&bal, Balance{
					AppID:   a.AppID,
					AssetID: a.AssetID,
				}).Preload("Asset").Error

				if err != nil {
					u.Errorln("get balance error", a.ID)
					return err
				}
				u.Debugln("get balance record", a.ID)

				data := GetBalanceInitData()

				data["InLocked"].Oper = "ADD"
				data["InLocked"].Value = a.Amount
				data["InLockedFee"].Oper = "ADD"
				data["InLockedFee"].Value = a.Fee

				err = ComputeBalance(tx, &bal, &data)
				if err != nil {
					u.Errorln("compute balance error", a.ID)
					return err
				}
				u.Debugln("compute balance error", a.ID)

				// save order
				order := new(Order)
				order.TotalAmount = a.TotalAmount
				order.Amount = a.Amount
				order.Fee = a.Fee
				order.AssetID = a.AssetID
				order.AppID = a.AppID
				order.JPOrderID = a.ID
				order.JPHash = a.Hash
				order.JPUUHash = a.UUHash
				order.Status = OrderStatusInit
				order.Type = OrderTypeDeposit
				order.Settled = false
				tx.Save(order)
				if err != nil {
					u.Errorf("save order error,", err, a.ID)
					return err
				}
			} else if a.Status == JPOrderStatusFailed { // case 2
				// DEPOSIT JPOrder NOT settled before
				// status: -> FAILED
				// balance: InLock -= 0, balance -= 0
				// do NOTHING

			} else if a.Status == JPOrderStatusPending { // case 3
				// DEPOSIT JPOrder NOT settled before
				// status: -> PENDING
				// balance: InLock += jporder.amount, balance += 0, InLockedFee += fee from asset same as case 1

			}
		} else if a.Settled {
			if a.Status == JPOrderStatusDone { // case 4
				// DEPOSIT JPOrder settled before
				// status: PENDING -> DONE
				// balance: InLock -= 0, balance += 0, InLockedFee -= 0
				// create ORDER, order.TotalAmount = jporder.amount, order.Fee = from Asset table, order.Amount = order.TotalAmount - order.Fee, same as case 1

			} else if a.Status == JPOrderStatusFailed { // case 5, symmetrical to case 3 & 1
				// DEPOSIT JPOrder settled before
				// status: PENDING -> FAILED
				// balance: InLock -= jporder.amount, balance += 0, InLockedFee -= fee

			} else if a.Status == JPOrderStatusPending { // case 6
				// DEPOSIT JPOrder settled before
				// status: PENDING -> PENDING
				// do NOTHING

			}
		}
	} else if a.Type == JPOrderTypeWithdraw { // map to cyborder's part
		if a.Status == JPOrderStatusDone {
			// status: -> DONE
			// balance: balance -= 0, outLocked -= TotalAmount, outLockedFee -= Fee

		} else if a.Status == JPOrderStatusFailed {
			// status: -> FAILED
			// balance: balance -= 0, outLocked -= 0, outLockedFee -= 0
			// create NEW jporder - status INIT, set it to order, move old jporder to order's FailedJPOrders

		} else if a.Status == JPOrderStatusPending {
			// status: -> PENDING
			// do NOTHING
		}
	}

	// return errors.New("test error for rollback")

	return nil
}
