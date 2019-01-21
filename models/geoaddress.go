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
