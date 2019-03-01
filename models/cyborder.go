package model

import (
	u "coding.net/bobxuyang/cy-gateway-BN/utils"
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
	// CybOrderTypeUR ...
	CybOrderTypeUR = "UR"
	// CybOrderTypeSweep ...
	CybOrderTypeSweep = "SWEEP"
	// CybOrderTypeFeeSettle ...
	CybOrderTypeFeeSettle = "FEESETTLE"
)

//CybOrder ...
type CybOrder struct {
	gorm.Model

	AssetID uint  `gorm:"not null" json:"assetID"` // 1 to n
	AppID   uint  `gorm:"not null" json:"appID"`   // 1 to n
	Asset   Asset `gorm:"ForeignKey:AssetId" json:"asset"`
	// Accounting      Accounting `gorm:"foreignkey:AccountingRefer" json:"accounting"` // 1 to 1
	// AccountingRefer uint       `json:"accountingRefer"`

	From   string       `gorm:"type:varchar(128)" json:"from"`
	To     string       `gorm:"type:varchar(128)" json:"to"`
	CybFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"cybFee"` // cybex blockchain fee
	// FeeAssetID uint

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`      //
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`         // fee in Asset

	Hash         string `gorm:"unique;index;type:varchar(128);default:null" json:"hash"`   // block:index
	UUHash       string `gorm:"unique;index;type:varchar(256);default:null" json:"uuhash"` // used for signature
	Sig          string `gorm:"unique;index;type:varchar(256);default:null" json:"sig"`    // used for signature
	Status       string `gorm:"type:varchar(32);not null" json:"status"`                   // INIT, HOLDING, PENDING, DONE, FAILED
	Type         string `gorm:"type:varchar(32);not null" json:"type"`                     // DEPOSIT, WITHDRAW, RECHARGE, SWEEP, FEESETTLE
	Memo         string `json:"Memo"`
	WithdrawAddr string `json:"withdrawAddr"`
	Settled      bool   `gorm:"not null;default:false" json:"settled"`   // if count amount to balance, then Settled = true
	Finalized    bool   `gorm:"not null;default:false" json:"finalized"` // if jporder was done or failed before
	EnterHook    bool   `gorm:"not null;default:false" json:"enterHook"` // set it to true if biz-logic need go-through after-save hook
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

// balance: balance += cyborder.Amount, inLocked -= cyborder.TotalAmount, inLockedFee -= cyborder.Fee, case 7
func (a *CybOrder) computeInLocked(tx *gorm.DB, oper string) error {
	var bal Balance
	err := tx.FirstOrCreate(&bal, Balance{AppID: a.AppID, AssetID: a.AssetID}).Error

	if err != nil {
		u.Errorln("get balance error", a.ID)
		return err
	}
	u.Debugln("get balance record", a.ID)

	moper := "SUB"
	if oper == "SUB" {
		moper = "ADD"
	}

	data := GetBalanceInitData()

	data["Balance"].Oper = moper
	data["Balance"].Value = a.Amount // deposit Amount to user, not TotalAmount
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

// balance: OutLock += cyborder.TotalAmount, OutLockedFee +=cyborder.Fee, balance -= cyborder.TotalAmount, case 1 & 3
// balance: OutLock -= cyborder.TotalAmount, OutLockedFee -=cyborder.Fee, balance += cyborder.TotalAmount, case 5
func (a *CybOrder) computeOutLocked(tx *gorm.DB, oper string) error {
	var bal Balance
	err := tx.FirstOrCreate(&bal, Balance{AppID: a.AppID, AssetID: a.AssetID}).Error

	if err != nil {
		u.Errorln("get balance error", a.ID)
		return err
	}
	u.Debugln("get balance record", a.ID)

	moper := "SUB"
	if oper == "SUB" {
		moper = "ADD"
	}

	data := GetBalanceInitData()

	data["Balance"].Oper = moper
	data["Balance"].Value = a.TotalAmount
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
// create ORDER, order.TotalAmount = cyborder.TotalAmount, order.Fee = cyborder.Fee, order.Amount = cyborder.Amount, case 1 & 4
func (a *CybOrder) CreateOrder(tx *gorm.DB) error {
	// save order
	order := new(Order)
	order.TotalAmount = a.TotalAmount
	order.Amount = a.Amount
	order.Fee = a.Fee
	order.AssetID = a.AssetID
	order.AppID = a.AppID
	order.CybOrderID = a.ID
	order.CybHash = a.Hash
	order.CybUUHash = a.UUHash
	order.Status = OrderStatusInit
	order.Type = OrderTypeWithdraw

	return tx.Save(order).Error
}

// Clone ...
func (a *CybOrder) Clone() *CybOrder {
	// save order
	order := new(CybOrder)

	order.AssetID = a.AssetID
	order.AppID = a.AppID

	order.To = a.To
	order.TotalAmount = a.TotalAmount
	order.Amount = a.Amount
	order.Fee = a.Fee
	order.Type = a.Type

	order.Status = CybOrderStatusInit
	order.Settled = false
	order.Finalized = false

	return order
}

//AfterSaveHook ... should be called manually
func (a CybOrder) AfterSaveHook(tx *gorm.DB) (err error) {
	if a.Finalized || !a.EnterHook {
		return nil
	}

	a.EnterHook = false
	err = tx.Save(a).Error
	if err != nil {
		u.Errorf("save jporder error,", err, a.ID)
		return err
	}

	if a.Type == CybOrderTypeWithdraw {
		if a.Settled == false {
			// case 1, 2, 4, 6 will executed only once

			// set CybOrder Settled = true and SAVE to DB
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

			if a.Status == CybOrderStatusDone { // case 1
				// DEPOSIT CybOrder NOT settled before
				// status: -> DONE
				// balance: OutLock += cyborder.TotalAmount, OutLockedFee +=cyborder.Fee, balance -= cyborder.TotalAmount, same as case 3
				// update ORDER

				// update order status to init
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{CybOrderID: a.ID}).First(&order).Error
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
				err = a.computeOutLocked(tx, "ADD")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}

			} else if a.Status == CybOrderStatusFailed { // case 2
				// DEPOSIT CybOrder NOT settled before
				// status: -> FAILED

				// update order status to failed
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{CybOrderID: a.ID}).First(&order).Error
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

			} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding { // case 3
				// DEPOSIT CybOrder NOT settled before
				// status: -> PENDING
				// balance: OutLock += cyborder.TotalAmount, OutLockedFee +=cyborder.Fee, balance -= cyborder.TotalAmount, same as case 1

				err = a.computeOutLocked(tx, "ADD")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}
			}
		} else if a.Settled {
			if a.Status == CybOrderStatusDone { // case 4
				// DEPOSIT CybOrder settled before
				// status: PENDING -> DONE
				// same as case 1

				// update order status to init
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{CybOrderID: a.ID}).First(&order).Error
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
			} else if a.Status == CybOrderStatusFailed { // case 5, symmetrical to case 3 & 1
				// DEPOSIT CybOrder settled before
				// status: PENDING -> FAILED
				// balance: OutLock -= cyborder.TotalAmount, OutLockedFee -=cyborder.Fee, balance += cyborder.TotalAmount

				// update order status to failed
				order := Order{}
				err = tx.Model(&Order{}).Where(&Order{CybOrderID: a.ID}).First(&order).Error
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

				err = a.computeOutLocked(tx, "SUB")
				if err != nil {
					u.Errorf("compute balance error,", err, a.ID)
					return err
				}
			} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding { // case 6
				// DEPOSIT CybOrder settled before
				// status: PENDING -> PENDING
				// do NOTHING

			}
		}
	} else if a.Type == CybOrderTypeDeposit { // map to jporder's part
		if a.Status == CybOrderStatusDone { // case 7
			// status: -> DONE
			// balance: balance += cyborder.Amount, inLocked -= cyborder.TotalAmount, inLockedFee -= cyborder.Fee

			err = a.computeInLocked(tx, "SUB")
			if err != nil {
				u.Errorf("compute balance error,", err, a.ID)
				return err
			}
		} else if a.Status == CybOrderStatusFailed { // case 8
			// status: -> FAILED
			// balance: InLock -= 0, balance -= 0, inLockedFee -=0
			// create NEW cyborder, set it to order, move old cyborder to order's FailedCybOrders

			b := a.Clone()
			if err != nil {
				u.Errorf("clone cyborder error,", err, a.ID)
				return err
			}

			err = tx.Save(b).Error
			if err != nil {
				u.Errorf("create cyborder error,", err, a.ID)
				return err
			}

			// TODO ...
		} else if a.Status == CybOrderStatusPending || a.Status == CybOrderStatusInit || a.Status == CybOrderStatusHolding { // case 9
			// status: -> PENDING
			// do NOTHING
		}
	}

	if a.Status == CybOrderStatusDone || a.Status == CybOrderStatusFailed {
		a.Finalized = true
		err := tx.Save(a).Error
		if err != nil {
			u.Errorf("set cyborder's Finalized to true error,", err, a.ID)
			return err
		}
	}

	return nil
}
