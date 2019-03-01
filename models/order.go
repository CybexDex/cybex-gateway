package model

import (
	u "coding.net/bobxuyang/cy-gateway-BN/utils"
	utype "coding.net/bobxuyang/cy-gateway-BN/utils/types"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

const (
	// OrderStatusPreInit ...
	OrderStatusPreInit = "PREINIT"
	// OrderStatusInit ...
	OrderStatusInit = "INIT"
	// OrderStatusProcessing ...
	OrderStatusProcessing = "PROCESSING"
	// OrderStatusDone ...
	OrderStatusDone = "DONE"
	// OrderStatusTerminated ...
	OrderStatusTerminated = "TERMINATED"
	// OrderStatusFailed ...
	OrderStatusFailed = "FAILED"
	// OrderStatusWaiting ...
	OrderStatusWaiting = "WAITING"
	// OrderTypeDeposit ...
	OrderTypeDeposit = "DEPOSIT"
	// OrderTypeWithdraw ...
	OrderTypeWithdraw = "WITHDRAW"
)

// RecordsQuery ...
type RecordsQuery struct {
	FundType string `schema:"fundType"`
	Offset   string `schema:"offset"`
	Size     string `schema:"size"`
	Asset    string `schema:"asset"`
	AppID    uint
}

// RecordsOut ...
type RecordsOut struct {
	*Order
	JPOrderID  utype.Omit `json:"jPOrderID,omitempty"`
	CybOrderID utype.Omit `json:"cybOrderID,omitempty"`
	Asset      string     `json:"asset"`
}

//Order ...
type Order struct {
	gorm.Model

	JPOrderID  uint      `json:"jPOrderID"`
	JPOrder    *JPOrder  `json:"jpOrder"`
	CybOrderID uint      `json:"cybOrderID"`
	CybOrder   *CybOrder `json:"cybOrder"`

	AssetID uint   `gorm:"not null" json:"-"` // 1 to n
	Asset   *Asset `json:"asset"`
	AppID   uint   `gorm:"not null" json:"-"` // 1 to n

	FailedJPOrders  pq.Int64Array `gorm:"type:integer[]" json:"-"`
	FailedCybOrders pq.Int64Array `gorm:"type:integer[]" json:"-"`

	JPHash    string `gorm:"index;type:varchar(128)" json:"-"` //
	JPUUHash  string `gorm:"index;type:varchar(256)" json:"-"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	CybHash   string `gorm:"index;type:varchar(128)" json:"-"` //
	CybUUHash string `gorm:"index;type:varchar(256)" json:"-"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)

	TotalAmount *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"totalAmount"` // totalAmount = amount + fee
	Amount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"amount"`      //
	Fee         *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`         // fee in Asset

	Status string `gorm:"type:varchar(32);not null" json:"status"` // INIT, PROCESSING, DONE, TERMINATED
	Type   string `gorm:"type:varchar(32);not null" json:"type"`   // DEPOSIT, WITHDRAW

	Settled   bool `gorm:"not null;default:false" json:"settled"`   // if the third phase's order is created
	Finalized bool `gorm:"not null;default:false" json:"finalized"` // if order was done or terminated before
	EnterHook bool `gorm:"not null;default:false" json:"enterHook"` // set it to true if biz-logic need go-through after-save hook
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

func (a *Order) createCybOrder(tx *gorm.DB) (*CybOrder, error) {
	return nil, nil
}

func (a *Order) createCYBOrder(tx *gorm.DB) (*CybOrder, error) {
	// save order
	order := new(CybOrder)

	order.AssetID = a.AssetID
	order.AppID = a.AppID

	app := App{}
	err := tx.First(&app, a.AppID).Error
	if err != nil {
		u.Errorf("get app error,", err, a.ID)
		return nil, err
	}
	order.To = app.CybAccount

	order.TotalAmount = a.TotalAmount
	order.Amount = a.Amount
	order.Fee = a.Fee

	order.Type = a.Type
	order.Status = CybOrderStatusInit

	order.Settled = false
	order.Finalized = false

	return order, tx.Save(order).Error
}

func (a *Order) createJPOrder(tx *gorm.DB) (*JPOrder, error) {
	// save order
	order := new(JPOrder)

	order.AssetID = a.AssetID
	order.AppID = a.AppID

	app := App{}
	err := tx.First(&app, a.AppID).Error
	if err != nil {
		u.Errorf("get app error,", err, a.ID)
		return nil, err
	}
	order.JadepoolID = app.JadepoolID

	o := CybOrder{}
	err = tx.First(&o, a.CybOrderID).Error
	if err != nil {
		u.Errorf("get cyborder error,", err, a.ID)
		return nil, err
	}
	order.To = o.WithdrawAddr

	order.TotalAmount = a.TotalAmount
	order.Amount = a.Amount
	order.Fee = a.Fee
	order.Confirmations = 0

	order.Type = a.Type
	order.Status = JPOrderStatusInit

	order.Resend = false
	order.Settled = false
	order.Finalized = false

	return order, tx.Save(order).Error
}

// CreateNext ...
func (a *Order) CreateNext(tx *gorm.DB) (err error) {
	if a.Status == OrderStatusDone && a.Type == OrderTypeDeposit {
		// create cyborder
		cybOrder, err := a.createCYBOrder(tx)
		if err != nil {
			u.Errorf("save cyborder error,", err, a.ID)
			return err
		}
		a.CybOrderID = cybOrder.ID
	} else if a.Status == OrderStatusDone && a.Type == OrderTypeWithdraw {
		// create jporder
		jpOrder, err := a.createJPOrder(tx)
		if err != nil {
			u.Errorf("save jporder error,", err, a.ID)
			return err
		}
		a.JPOrderID = jpOrder.ID
	}
	return nil
}

//AfterSaveHook ... should be called manually
func (a *Order) AfterSaveHook(tx *gorm.DB) (err error) {
	if a.Finalized || !a.EnterHook {
		return nil
	}

	a.EnterHook = false
	err = tx.Save(a).Error
	if err != nil {
		u.Errorf("save order error,", err, a.ID)
		return err
	}

	if a.Status == OrderStatusDone && a.Type == OrderTypeDeposit {
		// create cyborder
		cybOrder, err := a.createCYBOrder(tx)
		if err != nil {
			u.Errorf("save cyborder error,", err, a.ID)
			return err
		}
		a.CybOrderID = cybOrder.ID
		err = tx.Save(a).Error
		if err != nil {
			u.Errorf("save order error,", err, a.ID)
			return err
		}
	} else if a.Status == OrderStatusDone && a.Type == OrderTypeWithdraw {
		// create jporder
		jpOrder, err := a.createJPOrder(tx)
		if err != nil {
			u.Errorf("save jporder error,", err, a.ID)
			return err
		}
		a.JPOrderID = jpOrder.ID
		err = tx.Save(a).Error
		if err != nil {
			u.Errorf("save order error,", err, a.ID)
			return err
		}
	}

	if a.Status == OrderStatusTerminated && a.Type == OrderTypeDeposit {
		// do NOTHING
	} else if a.Status == OrderStatusTerminated && a.Type == OrderTypeWithdraw {
		// do NOTHING
	}

	if a.Status == OrderStatusDone || a.Status == OrderStatusTerminated {
		a.Finalized = true
		err := tx.Save(a).Error
		if err != nil {
			u.Errorf("set order's Finalized to true error,", err, a.ID)
			return err
		}
	}

	return nil
}
