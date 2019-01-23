package model

import "github.com/jinzhu/gorm"

//Address ...
type Address struct {
	gorm.Model

	ExOrders []ExOrder `json:"exOrders"`
	Orders   []Order   `json:"Orders"`

	AppID   uint `json:"appID"`
	AssetID uint `json:"assetID"`

	Address string `gorm:"index;type:varchar(128);not null" json:"address"`
	Status  string `gorm:"type:varchar(32);not null" json:"status"` // INIT, NORMAL, ABNORMAL
	// memo      string `gorm:"type:varchar(64)" json:"memo"`
	// UUAddress string `gorm:"type:varchar(255);not null" json:"uuaddres"` // = BLOCKCHAIN + ADDRESS
}

//FetchAll ...
func (add *Address) FetchAll() (res []*Address, err error) {
	err = GetDB().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (add *Address) Fetch(p Page) (res []*Address, err error) {
	err = GetDB().Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (add *Address) GetByID(id uint) (*Address, error) {
	a := Address{}
	err := GetDB().First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//Update ...
func (add *Address) Update(id uint, v *Address) error {
	return GetDB().Model(Address{}).Where("ID=?", id).UpdateColumns(v).Error
}

//Create ...
func (add *Address) Create(a *Address) (err error) {
	return GetDB().Create(&a).Error
}

//DeleteByID ...
func (add *Address) DeleteByID(id uint) (err error) {
	return GetDB().Where("ID=?", id).Delete(&Address{}).Error
}

//Delete ...
func (add *Address) Delete(a *Address) (err error) {
	return GetDB().Delete(&a).Error
}
