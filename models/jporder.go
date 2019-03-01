package model

import (
	u "coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

const (
	// JPOrderStatusInit ...
	JPOrderStatusInit = "INIT"
	// JPOrderStatusHolding ...
	JPOrderStatusHolding = "HOLDING"
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

	AssetID         uint `gorm:"not null" json:"assetID"`    // n to 1
	JadepoolID      uint `gorm:"not null" json:"jadepoolID"` // n to 1
	AppID           uint `gorm:"not null" json:"appID"`      // n to 1
	JadepoolOrderID uint `json:"jadepoolOrderID"`            // n to 1

	Index         int          `json:"index"`                                           //
	Hash          string       `gorm:"index;type:varchar(128)" json:"hash"`             //
	UUHash        string       `gorm:"index;type:varchar(256)" json:"uuhash"`           // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	From          string       `gorm:"type:varchar(128)" json:"from"`                   //
	To            string       `gorm:"type:varchar(128)" json:"to"`                     //
	TotalAmount   *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount        *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`      //
	Fee           *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`         // fee in Asset
	Confirmations int          `json:"confirmations"`                                   //
	Resend        bool         `gorm:"not null;default:false" json:"resend"`            //
	Status        string       `gorm:"type:varchar(32);not null" json:"status"`         // INIT, HOLDING, PENDING, DONE, FAILED
	Type          string       `gorm:"type:varchar(32);not null" json:"type"`           // DEPOSIT, WITHDRAW

	Settled   bool `gorm:"not null;default:false" json:"settled"`   // if count amount to balance, then Settled = true
	Finalized bool `gorm:"not null;default:false" json:"finalized"` // if jporder was done or failed before
	EnterHook bool `gorm:"not null;default:false" json:"enterHook"` // set it to true if biz-logic need go-through after-save hook
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

// computeInLocked ...
// balance: InLock += jporder.TotalAmount, balance += 0, InLockedFee += jporder.Fee, case 1 & 3
// balance: InLock -= jporder.TotalAmount, balance += 0, InLockedFee -= jporder.Fee, case 5
func (a *JPOrder) computeInLocked(tx *gorm.DB, oper string) error {
	var bal Balance
	err := tx.FirstOrCreate(&bal, Balance{AppID: a.AppID, AssetID: a.AssetID}).Error

	if err != nil {
		u.Errorln("get balance error", a.ID)
		return err
	}
	u.Debugln("get balance record", a.ID)

	data := GetBalanceInitData()

	data["InLocked"].Oper = oper
	data["InLocked"].Value = a.TotalAmount
	data["InLockedFee"].Oper = oper
	data["InLockedFee"].Value = a.Fee

	err = ComputeBalance(tx, &bal, &data)
	if err != nil {
		u.Errorln("compute balance error", err, a.ID)
		return err
	}
	u.Debugln("compute balance", a.ID)

	return nil
}

//computeOutLocked ...
// balance: balance -= 0, outLocked -= TotalAmount, outLockedFee -= Fee, case 7
func (a *JPOrder) computeOutLocked(tx *gorm.DB, oper string) error {
	var bal Balance
	err := tx.FirstOrCreate(&bal, Balance{AppID: a.AppID, AssetID: a.AssetID}).Error

	if err != nil {
		u.Errorln("get balance error", a.ID)
		return err
	}
	u.Debugln("get balance record", a.ID)

	data := GetBalanceInitData()

	data["OutLocked"].Oper = oper
	data["OutLocked"].Value = a.TotalAmount
	data["OutLockedFee"].Oper = oper
	data["OutLockedFee"].Value = a.Fee

	err = ComputeBalance(tx, &bal, &data)
	if err != nil {
		u.Errorln("compute balance error", err, a.ID)
		return err
	}
	u.Debugln("compute balance", a.ID)

	return nil
}

// CreateOrder ...
// create ORDER, order.TotalAmount = jporder.TotalAmount, order.Fee = jporder.Fee, order.Amount = jporder.Amount, case 1 & 4
func (a *JPOrder) CreateOrder(tx *gorm.DB) error {
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
	order.Status = OrderStatusPreInit
	order.Type = OrderTypeDeposit
	order.Settled = false
	order.Finalized = false

	return tx.Save(order).Error
}

// Clone ...
func (a *JPOrder) Clone() *JPOrder {
	// save order
	order := new(JPOrder)

	order.AssetID = a.AssetID
	order.AppID = a.AppID
	order.JadepoolID = a.JadepoolID
	order.JadepoolOrderID = a.JadepoolOrderID

	order.Index = a.Index
	order.To = a.To
	order.TotalAmount = a.TotalAmount
	order.Amount = a.Amount
	order.Fee = a.Fee
	order.Confirmations = 0
	order.Type = a.Type

	order.Status = JPOrderStatusInit
	order.Resend = false
	order.Settled = false
	order.Finalized = false

	return order
}

//AfterSaveHook ... should be called manually
func (a *JPOrder) AfterSaveHook(tx *gorm.DB) (err error) {
	if a.Finalized || !a.EnterHook {
		return nil
	}

	a.EnterHook = false
	err = tx.Save(a).Error
	if err != nil {
		u.Errorf("save jporder error,", err, a.ID)
		return err
	}

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

			// create order
			err = a.CreateOrder(tx)
			if err != nil {
				u.Errorf("save order error,", err, a.ID)
				return err
			}

			if a.Status == JPOrderStatusDone { // case 1
				// DEPOSIT JPOrder NOT settled before
				// status: -> DONE
				// balance: InLock += jporder.TotalAmount, balance += 0, InLockedFee += jporder.Fee, same as case 3
				// order.TotalAmount = jporder.TotalAmount, order.Fee = jporder.Fee, order.Amount = jporder.Amount

				// update order status to init
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{JPOrderID: a.ID}).First(&order).Error
				if err != nil {
					u.Errorf("get order error,", err, a.ID)
					return err
				}
				order.Status = OrderStatusInit
				err = tx.Save(order).Error
				if err != nil {
					u.Errorf("save order error,", err, a.ID)
					return err
				}

				// update balance record
				err = a.computeInLocked(tx, "ADD")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}
			} else if a.Status == JPOrderStatusFailed { // case 2
				// DEPOSIT JPOrder NOT settled before
				// status: -> FAILED
				// update order status to failed

				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{JPOrderID: a.ID}).First(&order).Error
				if err != nil {
					u.Errorf("get order error,", err, a.ID)
					return err
				}
				order.Status = OrderStatusFailed
				err = tx.Save(order).Error
				if err != nil {
					u.Errorf("save order error,", err, a.ID)
					return err
				}
			} else if a.Status == JPOrderStatusPending { // case 3
				// DEPOSIT JPOrder NOT settled before
				// status: -> PENDING
				// balance: InLock += jporder.TotalAmount, balance += 0, InLockedFee += jporder.Fee from asset same as case 1

				err = a.computeInLocked(tx, "ADD")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}
			}
		} else if a.Settled {
			if a.Status == JPOrderStatusDone { // case 4
				// DEPOSIT JPOrder settled before
				// status: PENDING -> DONE
				// same as case 1

				// update order status to init
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{JPOrderID: a.ID}).First(&order).Error
				if err != nil {
					u.Errorf("get order error,", err, a.ID)
					return err
				}
				order.Status = OrderStatusInit
				err = tx.Save(order).Error
				if err != nil {
					u.Errorf("save order error,", err, a.ID)
					return err
				}
			} else if a.Status == JPOrderStatusFailed { // case 5, symmetrical to case 3 & 1
				// DEPOSIT JPOrder settled before
				// status: PENDING -> FAILED
				// balance: InLock -= jporder.TotalAmount, balance += 0, InLockedFee -= jporder.Fee

				// update order status to failed
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{JPOrderID: a.ID}).First(&order).Error
				if err != nil {
					u.Errorf("get order error,", err, a.ID)
					return err
				}
				order.Status = OrderStatusFailed
				err = tx.Save(order).Error
				if err != nil {
					u.Errorf("save order error,", err, a.ID)
					return err
				}

				err = a.computeInLocked(tx, "SUB")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}
			} else if a.Status == JPOrderStatusPending { // case 6
				// DEPOSIT JPOrder settled before
				// status: PENDING -> PENDING
				// do NOTHING

			}
		}
	} else if a.Type == JPOrderTypeWithdraw { // map to cyborder's part
		if a.Status == JPOrderStatusDone { // case 7
			// status: -> DONE
			// balance: balance -= 0, outLocked -= TotalAmount, outLockedFee -= Fee

			err = a.computeOutLocked(tx, "SUB")
			if err != nil {
				u.Errorf("compute balance error,", err, a.ID)
				return err
			}
		} else if a.Status == JPOrderStatusFailed { // case 8
			// status: -> FAILED
			// balance: balance -= 0, outLocked -= 0, outLockedFee -= 0
			// create NEW jporder - status INIT, set it to order, move old jporder to order's FailedJPOrders

			if a.Resend {
				b := a.Clone()
				if err != nil {
					u.Errorf("clone jporder error,", err, a.ID)
					return err
				}

				err = tx.Save(b).Error
				if err != nil {
					u.Errorf("create jporder error,", err, a.ID)
					return err
				}

				// set it to order, move old jporder to order's FailedJPOrders
				order := Order{}
				err = tx.Where(&Order{JPOrderID: a.ID}).First(&order).Error
				if err != nil {
					u.Errorf("find order error: %v, jporder id: %d", err, a.ID)
					return err
				}
				order.JPOrderID = b.ID
				order.FailedJPOrders = append(order.FailedJPOrders, int64(a.ID))
				err = tx.Save(order).Error
				if err != nil {
					u.Errorf("save order error: %v", err)
					return err
				}
			}
		} else if a.Status == JPOrderStatusPending { // case 9
			// status: -> PENDING
			// do NOTHING
		}
	}

	if a.Status == JPOrderStatusDone || a.Status == JPOrderStatusFailed {
		a.Finalized = true
		err := tx.Save(a).Error
		if err != nil {
			u.Errorf("set jporder's Finalized to true error,", err, a.ID)
			return err
		}
	}

	return nil
}
