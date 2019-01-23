package model

import (
	"github.com/jinzhu/gorm"
)

//Account ...
type Account struct {
	gorm.Model

	Apps   []App   `gorm:"many2many:account_apps;" json:"apps"` // n to n, ???
	Events []Event `json:"events"`                              // people to event list

	CompanyID uint `json:"companyID"`

	Name            string `gorm:"unique;index;type:varchar(64);not null" json:"name"`
	LastName        string `gorm:"type:varchar(32)" json:"lastName"`
	FirstName       string `gorm:"type:varchar(32)" json:"firstName"`
	Email           string `gorm:"unique;index;type:varchar(64);not null" json:"email"`
	EmailVerified   bool   `gorm:"default:false" json:"emailVerified"`
	EmailEnable     bool   `gorm:"default:false" json:"emailEnabled"`
	Phone           string `gorm:"type:varchar(32)" json:"phone"`
	PhoneVerified   bool   `gorm:"default:false" json:"phoneVerified"`
	PhoneEnabled    bool   `gorm:"default:false" json:"phoneEnabled"`
	AuthKey         string `gorm:"type:varchar(64)" json:"authKey"`
	AuthKeyVerified bool   `gorm:"default:false" json:"authKeyVerified"`
	AuthKeyEnabled  bool   `gorm:"default:false" json:"authKeyEnabled"`
	PasswordHash    string `gorm:"type:varchar(512);not null" json:"passwordHash"`
	Password        string `gorm:"-" json:"password"`
	Status          string `gorm:"type:varchar(32);default:'INIT'" json:"status"`    // INIT, NORMAL, ABNORMAL
	Type            string `gorm:"type:varchar(32);default:'SAAS_USER'" json:"type"` // ADMIN, SUPER_ADMIN, SAAS_USER
	Disable         bool   `gorm:"default:false" json:"disable"`
}
