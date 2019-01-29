package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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

	FailedJPOrders  pq.Int64Array `gorm:"type:integer[]" json:"failedJPOrders"`
	FailedCybOrders pq.Int64Array `gorm:"type:integer[]" json:"failedCybOrders"`

	JPHash   string `gorm:"unique;index;type:varchar(128)" json:"jpHash"`
	JPUUHash string `gorm:"unique;index;type:varchar(256)" json:"jpUUHash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	CybHash   string `gorm:"unique;index;type:varchar(128)" json:"cybHash"`
	CybUUHash string `gorm:"unique;index;type:varchar(256)" json:"cybUUHash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"` // fee in Asset

	Status    string `gorm:"type:varchar(32);not null" json:"status"` // INIT, PROCESSING, DONE, TERMINATED
	Type      string `gorm:"type:varchar(32);not null" json:"type"`   // DEPOSIT, WITHDRAW
	Settled   bool   `gorm:"not null;default:false" json:"settled"`   // whether the third phase's order is created
	Finalized bool   `gorm:"not null;default:false" json:"finalized"` // if order was done or failed before
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

//AfterSave ... will be called each time after CREATE / SAVE / UPDATE
func (a Order) AfterSave(tx *gorm.DB) (err error) {
	if a.Settled == false {
		// set order's settled = true and SAVE to DB

		//保证只做一次

		if a.Status == OrderStatusDone && a.Type == OrderTypeDeposit {
			// create cyborder
		} else if a.Status == OrderStatusDone && a.Type == OrderTypeWithdraw {
			// create jporder
		}

		if a.Status == OrderStatusTerminated && a.Type == OrderTypeDeposit {
			// inlock -= amount

		} else if a.Status == OrderStatusTerminated && a.Type == OrderTypeWithdraw {
			// outlock -= amount, balance += amount
		}

	}

	// return errors.New("test error for rollback")

	return nil
}
