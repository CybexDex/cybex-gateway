package model

import "github.com/jinzhu/gorm"

//Company ...
type Company struct {
	gorm.Model

	Accounts []Account `gorm:"ForeignKey:CompanyID" json:"accounts"` // 1 to n
	Apps     []App     `json:"apps"`                                 // 1 to n

	CompanyAddress      GeoAddress `gorm:"foreignkey:CompanyAddressRefer" json:"companyAddress"` // 1 to 1
	CompanyAddressRefer uint       `json:"companyAddressRefer"`
	ContactAddress      GeoAddress `gorm:"foreignkey:ContactAddressRefer" json:"contactAddress"` // 1 to 1
	ContactAddressRefer uint       `json:"contactAddressRefer"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"type:varchar(32);not null" json:"status"` // INIT, NORMAL, UN_PAIED, ABNORMAL, DELETED
	Type        string `gorm:"type:varchar(32)" json:"type"`            // ??
}
