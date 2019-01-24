package model

import "github.com/jinzhu/gorm"

//GeoAddress ...
type GeoAddress struct {
	gorm.Model

	Address   string `gorm:"type:text;not null" json:"address"`
	Zipcode   string `gorm:"type:varchar(32);not null" json:"zipcode"`
	LastName  string `gorm:"type:varchar(32)" json:"lastName"`
	FirstName string `gorm:"type:varchar(32)" json:"firstName"`
	Phone     string `gorm:"type:varchar(32)" json:"phone"`
}

//UpdateColumns ...
func (a *GeoAddress) UpdateColumns(b *GeoAddress) error {
	return GetDB().Model(GeoAddress{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *GeoAddress) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *GeoAddress) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *GeoAddress) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
