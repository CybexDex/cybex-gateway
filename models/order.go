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

// GORMMOVE
type GORMMove struct {
	ID utype.Omit `json:"ID,omitempty"`
	// CreatedAt utype.Omit `json:"CreatedAt,omitempty"`
	UpdatedAt utype.Omit `json:"UpdatedAt,omitempty"`
	DeletedAt utype.Omit `json:"DeletedAt,omitempty"`
}

// RecordsOut ...
type RecordsOut struct {
	*Order
	*GORMMove
	CybexName   string     `json:"cybexName"`
	OutAddr     string     `json:"outAddr"`
	Confirms    string     `json:"confirms"`
	JPOrderID   utype.Omit `json:"jPOrderID,omitempty"`
	CybOrderID  utype.Omit `json:"cybOrderID,omitempty"`
	AssetID     utype.Omit `json:"assetID,omitempty"`
	AppID       utype.Omit `json:"appID,omitempty"`
	Asset       string     `json:"asset"`
	JPOrder     utype.Omit `json:"jpOrder,omitempty"`
	CybOrder    utype.Omit `json:"cybOrder,omitempty"`
	App         utype.Omit `json:"app,omitempty"`
	OutHash     string     `json:"outHash"`
	CybHash     string     `json:"cybHash"`
	TotalAmount string     `json:"totalAmount"` // totalAmount = amount + fee
	Amount      string     `json:"amount"`      //
	Fee         string     `json:"fee"`
	Status      string     `json:"status"`
}

//Order ...
type Order struct {
	gorm.Model

	JPOrderID  uint      `json:"jPOrderID"`
	JPOrder    *JPOrder  `json:"jpOrder"`
	CybOrderID uint      `json:"cybOrderID"`
	CybOrder   *CybOrder `json:"cybOrder"`

	AssetID         uint          `gorm:"not null" json:"assetID"` // 1 to n
	Asset           *Asset        `json:"asset"`
	AppID           uint          `gorm:"not null" json:"appID"` // 1 to n
	App             *App          `json:"app"`
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
		err := tx.Save(a).Error
		if err != nil {
			u.Errorf("set order's Finalized to true error,", err, a.ID)
			return err
		}
	}

	return nil
}
