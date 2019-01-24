package model

import (
	"github.com/jinzhu/gorm"
)

//CybToken ...
type CybToken struct {
	gorm.Model
	CybAccount string `gorm:"type:varchar(255)" json:"accountName"` // use for GATEWAY mode
	Signer     string `gorm:"type:varchar(255)" json:"signer"`      // user's Signer
	Expiration uint   `json:"expiration"`                           //timestamp seconds
}

//UpdateColumns ...
func (a *CybToken) UpdateColumns(b *CybToken) error {
	return GetDB().Model(CybToken{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *CybToken) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *CybToken) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *CybToken) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
