package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

// below FEE is charged by SAAS(gateway), not by blockchain
//
// for DEPOSIT order
// total-amount 	amount		fee
// 100				99			1
// ---------- phase I
// balance.inLocked += 100, balance.balance += 0
// order.TotalAmount = jporder.amount(100) --- when create order
// order.Fee(from asset table) = 1, order.Amount = order.TotalAmount(100) - order.Fee(1) = 99
// balance.inLockedFee += order.Fee(1)
// ---------- phase II
// ---------- phase III
// status DONE:  balance.balance = balance.balance + order.Amount(99), balance.inLocked -= order.TotalAmount(100), balance.inLockedFee -= order.Fee(1)
// status FAILED: balance.balance += 0, balance.inLocked -= 0, balance.inLockedFee -= 0, create new cyborder
//
// for WITHDRAW order
// total-amount 	amount		fee
// 100				99			1
// ---------- phase I
// balance.outLocked += total-amount(100), balance.balance -= total-amount(100)
// order.TotalAmount = cyborder.amount(100) --- when create order
// ---------- phase II
// order.Fee(from asset table) = 1, order.Amount = order.TotalAmount(100) - order.Fee(1) = 99
// balance.outLockedFee += fee(1)
// ---------- phase III
// status DONE: balance.balance -= 0, balance.outLocked -= order.TotalAmount(100), balance.outLockedFee -= order.Fee(1)
// status FAILED: balance.balance -= 0, balance.outLocked -= 0, balance.outLockedFee -= 0, create new jporder
//

//Balance ...
type Balance struct {
	gorm.Model

	AppID uint `json:"appID"`

	AssetID uint  `json:"assetID"`
	Asset   Asset `gorm:"ForeignKey:AssetID" json:"asset"`

	Balance   *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"balance"`
	InLocked  *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"inLocked"`  // after order DONE, add to balance
	OutLocked *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"outLocked"` // deduct from balance and create WITHDRAWED ORDER

	InLockedFee  *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"inLockedFee"`  // after order DONE, add to balance
	OutLockedFee *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"outLockedFee"` // deduct from balance and create WITHDRAWED ORDER
}

//UpdateColumns ...
func (a *Balance) UpdateColumns(b *Balance) error {
	return GetDB().Model(Asset{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Balance) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Balance) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Balance) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
