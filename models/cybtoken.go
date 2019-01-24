package model

import (
	"github.com/jinzhu/gorm"
)

//CybToken ...
type CybToken struct {
	gorm.Model
	CybAccount string `gorm:"unique;type:varchar(255)" json:"accountName"` // use for GATEWAY mode
	Signer     string `gorm:"type:varchar(255)" json:"signer"`             // user's Signer
	Expiration uint   `json:"expiration"`                                  //timestamp seconds
}

//UpdateColumns ...
func (a *CybToken) UpdateColumns(b *CybToken) error {
	return GetDB().Model(CybToken{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *CybToken) Create() (err error) {
	return GetDB().Create(&a).Error
}

//SaveUniqueBy ...
func (a *CybToken) SaveUniqueBy(uniq CybToken) (err error) {
	return GetDB().Where(uniq).Assign(*a).FirstOrCreate(&a).Error
}

//Delete ...
func (a *CybToken) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
