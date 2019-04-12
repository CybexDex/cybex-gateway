package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
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
	User       string `json:"user"`
	BNOrderID  string `gorm:"index;" json:"bnOrderID"` // n to 1

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

	Current       string `json:"-"`
	CurrentState  string `json:"-"`
	CurrentReason string `json:"-"`
}

// JPOrderFind ...
func JPOrderFind(j *JPOrder) (res []*JPOrder, err error) {
	err = db.Where(j).Find(&res).Error
	return res, err
}

// Update ...
func (j *JPOrder) Update(i *JPOrder) error {
	return db.Model(JPOrder{}).Where("ID=?", j.ID).UpdateColumns(i).Error
}

// Save ...
func (j *JPOrder) Save() error {
	return db.Save(j).Error
}

// SetCurrent ...
func (j *JPOrder) SetCurrent(current string, state string, reason string) {
	j.Current = current
	j.CurrentState = state
	j.CurrentReason = reason
}

// JPOrderCreate ...
func JPOrderCreate(j *JPOrder) error {
	err := db.Create(j).Error
	if err != nil {
		return err
	}
	return nil
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
