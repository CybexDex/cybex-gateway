package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
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

// JPOrderCreate ...
func JPOrderCreate(j *JPOrder) error {
	err := db.Create(j).Error
	if err != nil {
		return err
	}
	return nil
}
