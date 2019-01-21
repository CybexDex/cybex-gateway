package model

import "github.com/jinzhu/gorm"

//Account ...
type Account struct {
	gorm.Model

	Apps   []App   `gorm:"many2many:account_apps;" json:"apps"` // n to n, ???
	Events []Event `json:"events"`                              // people to event list

	CompanyID uint `json:"companyID"`

	Name            string `gorm:"index;type:varchar(255);not null" json:"name"`
	LastName        string `gorm:"type:varchar(32)" json:"lastName"`
	FirstName       string `gorm:"type:varchar(32)" json:"firstName"`
	Email           string `gorm:"type:varchar(64)" json:"email"`
	EmailVerified   bool   `gorm:"default:false" json:"emailVerified"`
	EmailEnable     bool   `gorm:"default:false" json:"emailEnabled"`
	Phone           string `gorm:"type:varchar(32)" json:"phone"`
	PhoneVerified   bool   `gorm:"default:false" json:"phoneVerified"`
	PhoneEnabled    bool   `gorm:"default:false" json:"phoneEnabled"`
	AuthKey         string `gorm:"type:varchar(64)" json:"authKey"`
	AuthKeyVerified bool   `gorm:"default:false" json:"authKeyVerified"`
	AuthKeyEnabled  bool   `gorm:"default:false" json:"authKeyEnabled"`
	PasswordHash    string `gorm:"type:varchar(512);not null" json:"passwordHash"`
	Status          string `json:"status"` // INIT, NORMAL, ABNORMAL
	Type            string `json:"type"`   // ADMIN, SUPER_ADMIN, SAAS_USER
	Disable         bool   `gorm:"default:false" json:"disable"`
}
