package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Accounting ...
type Accounting struct {
	gorm.Model

	// AssetID       uint         `gorm:"assetID"` // asset pay for blockchain fee
	BlockchainFee *apd.Decimal `gorm:"type:numeric(30,10)" json:"blockchainFee"`
	InAmount      *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"inAmount"`
	OutAmount     *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"outAmount"`
	Fee           *apd.Decimal `gorm:"type:numeric(30,10);not null" json:"fee"`
	Settled       bool         `gorm:"default:false" json:"settled"`
}

//FetchAll ...
func (acc *Accounting) FetchAll() ([]*Accounting, error) {
	var res []*Accounting
	err := GetDB().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (acc *Accounting) Fetch(p Page) (res []*Accounting, err error) {
	err = GetDB().Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (acc *Accounting) GetByID(id uint) (*Accounting, error) {
	a := Accounting{}
	err := GetDB().First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//Update ...
func (acc *Accounting) Update(id uint, v *Account) error {
	return GetDB().Model(Accounting{}).Where("ID=?", id).UpdateColumns(v).Error
}

//Create ...
func (acc *Accounting) Create(a *Accounting) (err error) {
	return GetDB().Create(&a).Error
}

//DeleteByID ...
func (acc *Accounting) DeleteByID(id uint) (err error) {
	return GetDB().Where("ID=?", id).Delete(&Accounting{}).Error
}

//Delete ...
func (acc *Accounting) Delete(a *Accounting) (err error) {
	return GetDB().Delete(&a).Error
}
