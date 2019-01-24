package model

import (
	"github.com/jinzhu/gorm"
)

//Balance ...
type CybLogin struct {
	gorm.Model
	CybAccount string `gorm:"type:varchar(255)" json:"accountName"` // use for GATEWAY mode
	Signer     string `gorm:"type:varchar(255)" json:"signer"`      // user's Signer
	Expiration uint   `json:"expiration"`                           //timestamp seconds
}

func (acc *CybLogin) Create(a *CybLogin) (err error) {
	return nil
}
func (acc *CybLogin) Update(a *CybLogin) (err error) {
	return nil
}
