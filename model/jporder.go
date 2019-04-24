package model

import (
	"fmt"

	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

const (
	// JPOrderStatusInit ...
	JPOrderStatusInit = "INIT"
	// JPOrderStatusHolding ...
	JPOrderStatusHolding = "HOLDING"
	// JPOrderStatusProcessing ...
	JPOrderStatusProcessing = "PROCESSING"
	// JPOrderStatusPending ...
	JPOrderStatusPending = "PENDING"
	// JPOrderStatusDone ...
	JPOrderStatusDone = "DONE"
	// JPOrderStatusFailed ...
	JPOrderStatusFailed = "FAILED"
	//JPOrderStatusTerminate ...
	JPOrderStatusTerminate = "TERMINATE"
	// JPOrderTypeDeposit ...
	JPOrderTypeDeposit = "DEPOSIT"
	// JPOrderTypeWithdraw ...
	JPOrderTypeWithdraw = "WITHDRAW"
)

// JPOrder ...
type JPOrder struct {
	gorm.Model
	Asset      string `json:"asset"` // n to 1
	BlockChain string `json:"blockChain"`
	CybUser    string `json:"user"`
	OutAddr    string `json:"outAddr"`

	BNOrderID string `gorm:"index;" json:"bnOrderID"` // n to 1

	Index         int             `json:"index"`                                 //
	Hash          string          `gorm:"index;type:varchar(128)" json:"hash"`   //
	UUHash        string          `gorm:"index;type:varchar(256)" json:"uuhash"` // = BLOCKCHAINNAME + HASH + INDEX (if INDEX is null then ignore)
	From          string          `gorm:"type:varchar(128)" json:"from"`         //
	To            string          `gorm:"type:varchar(128)" json:"to"`           //
	Memo          string          `gorm:"type:varchar(256)" json:"memo"`
	TotalAmount   decimal.Decimal `json:"totalAmount" gorm:"type:numeric"`
	Amount        decimal.Decimal `json:"amount" gorm:"type:numeric"`
	Fee           decimal.Decimal `json:"fee" gorm:"type:numeric"`
	Confirmations int             `json:"confirmations"`                           //
	Resend        bool            `gorm:"not null;default:false" json:"resend"`    //
	Status        string          `gorm:"type:varchar(32);not null" json:"status"` // INIT, HOLDING, PENDING, DONE, FAILED
	Type          string          `gorm:"type:varchar(32);not null" json:"type"`   // DEPOSIT, WITHDRAW

	Link string `json:"link"`

	CYBHash  *string `gorm:"unique;index;type:varchar(128)" json:"-"`
	CYBHash2 string  `gorm:"index;type:varchar(128)" json:"-"`
	Sig      string  `json:"-"`
	Sig2     string  `json:"-"`

	Current       string `json:"-"`
	CurrentState  string `json:"-"`
	CurrentReason string `json:"-"`
}

// JPOrderFind ...
func JPOrderFind(j *JPOrder) (res []*JPOrder, err error) {
	err = db.Where(j).Find(&res).Error
	return res, err
}

// JPOrderRecord ...
func JPOrderRecord(user string, asset string, bizType string, size string, lastID string) (res []*JPOrder, count int, err error) {
	dbpre := db.Where(&JPOrder{
		CybUser: user,
		Type:    bizType,
		Asset:   asset,
	}).Where("id < ?", lastID).Order("id desc")
	err = dbpre.Limit(size).Find(&res).Error
	if err != nil {
		return nil, 0, err
	}
	xx := []*JPOrder{}
	err = dbpre.Find(&xx).Count(&count).Error
	return res, count, err
}

// Update ...
func (j *JPOrder) Update(i *JPOrder) error {
	return db.Model(JPOrder{}).Where("ID=?", j.ID).UpdateColumns(i).Error
}

// Save ...
func (j *JPOrder) Save() error {
	return db.Save(j).Error
}

// SetStatus ...
func (j *JPOrder) SetStatus(status string) {
	j.Status = status
	j.Log("SetStatus", fmt.Sprintln(status))
}

// SetCurrent ...
func (j *JPOrder) SetCurrent(current string, state string, reason string) {
	j.Current = current
	j.CurrentState = state
	j.CurrentReason = reason
	j.Log("SetCurrent", fmt.Sprintln(current, state, reason))
}

// Log ...
func (j *JPOrder) Log(event string, message string) {
	log1 := &OrderLog{
		OrderID: j.ID,
		Event:   event,
		Message: message,
	}
	db.Create(log1)
	log.Infof("order:%d,%s:%+v\n", log1.OrderID, "log", *log1)
}

// JPOrderCreate ...
func JPOrderCreate(j *JPOrder) error {
	err := db.Create(j).Error
	if err != nil {
		return err
	}
	return nil
}

// JPOrderUnBalanceInit ...
func JPOrderUnBalanceInit() {
	db.Model(JPOrder{}).Where(&JPOrder{
		Current:      "cybinner",
		CurrentState: "UNBALANCE",
	}).UpdateColumn(&JPOrder{
		CurrentState: "INIT",
	})
}

// HoldJPWithdrawOne ...
func HoldJPWithdrawOne() (*JPOrder, error) {
	var order1 JPOrder
	s := `update jp_orders 
	set current_state = 'PROCESSING' 
	where id = (
				select id 
				from jp_orders 
				where current_state = 'INIT' 
				and current = 'jp'
				and type = 'WITHDRAW'
				order by id
				limit 1
			)
	returning *`
	err := db.Raw(s).Scan(&order1).Error
	return &order1, err
}

//HoldCYBOrderOne ...
func HoldCYBOrderOne() (*JPOrder, error) {
	var order1 JPOrder
	s := `update jp_orders 
	set current_state = 'PROCESSING' 
	where id = (
				select id 
				from jp_orders 
				where current_state = 'INIT' 
				and current = 'cyborder'
				order by id
				limit 1
			)
	returning *`
	err := db.Raw(s).Scan(&order1).Error
	return &order1, err
}

//HoldCYBInnerOrderOne ...
func HoldCYBInnerOrderOne() (*JPOrder, error) {
	var order1 JPOrder
	s := `update jp_orders 
	set current_state = 'PROCESSING' 
	where id = (
				select id 
				from jp_orders 
				where current_state = 'INIT' 
				and type = 'WITHDRAW'
				and current = 'cybinner'
				order by id
				limit 1
			)
	returning *`
	err := db.Raw(s).Scan(&order1).Error
	return &order1, err
}

// HoldOrderOne ...
func HoldOrderOne() (*JPOrder, error) {
	var order1 JPOrder
	s := `update jp_orders 
	set current_state = 'PROCESSING' 
	where id = (
				select id 
				from jp_orders 
				where current_state = 'INIT' 
				and current = 'order'
				order by id
				limit 1
			)
	returning *`
	err := db.Raw(s).Scan(&order1).Error
	return &order1, err
}
