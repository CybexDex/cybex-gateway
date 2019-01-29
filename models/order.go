package model

import (
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
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

	JPHash    string `gorm:"unique;index;type:varchar(128)" json:"jpHash"`    //
	JPUUHash  string `gorm:"unique;index;type:varchar(256)" json:"jpUUHash"`  // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	CybHash   string `gorm:"unique;index;type:varchar(128)" json:"cybHash"`   //
	CybUUHash string `gorm:"unique;index;type:varchar(256)" json:"cybUUHash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`      //
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`         // fee in Asset

	Status    string `gorm:"type:varchar(32);not null" json:"status"` // INIT, PROCESSING, DONE, TERMINATED
	Type      string `gorm:"type:varchar(32);not null" json:"type"`   // DEPOSIT, WITHDRAW
	Settled   bool   `gorm:"not null;default:false" json:"settled"`   // if the third phase's order is created
	Finalized bool   `gorm:"not null;default:false" json:"finalized"` // if order was done or terminated before
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

func (a *Order) createCybOrder() (*CybOrder, error) {
	return nil, nil
}

func (a *Order) createJPOrder() (*JPOrder, error) {
	return nil, nil
}

//AfterSave1 ... will be called each time after CREATE / SAVE / UPDATE
func (a *Order) AfterSave1(tx *gorm.DB) (err error) {
	if a.Finalized {
		return nil
	}

	if a.Status == OrderStatusDone && a.Type == OrderTypeDeposit {
		// create cyborder

	} else if a.Status == OrderStatusDone && a.Type == OrderTypeWithdraw {
		// create jporder

	}

	if a.Status == OrderStatusTerminated && a.Type == OrderTypeDeposit {
		// do NOTHING

	} else if a.Status == OrderStatusTerminated && a.Type == OrderTypeWithdraw {
		// do NOTHING

	}

	if a.Status == OrderStatusDone || a.Status == OrderStatusTerminated {
		a.Finalized = true
		err := a.Save()
		if err != nil {
			u.Errorf("set jporder's Finalized to true error,", err, a.ID)
			return err
		}
	}

	return nil
}
