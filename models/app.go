package model

import "github.com/jinzhu/gorm"

//App ...
type App struct {
	gorm.Model

	Orders    []Order   `json:"orders"`    // 1 to n
	Balances  []Balance `json:"balances"`  // 1 to n, use for SASS mode
	Addresses []Address `json:"addresses"` // 1 to n

	CompanyID  uint `json:"companyID"`  // n to 1
	JadepoolID uint `json:"jadepoolID"` // n to 1

	Name        string `gorm:"index;type:varchar(255);not null" json:"name"`
	CybAccount  string `gorm:"type:varchar(255)" json:"cybAccount"` // use for GATEWAY mode
	Description string `gorm:"type:text" json:"description"`
	URL         string `gorm:"type:varchar(255)" json:"url"`
	Status      string `gorm:"type:varchar(32);not null" json:"status"` // INIT, NORMAL, UN_PAIED, ABNORMAL, DELETED
	Type        string `gorm:"type:varchar(32)" json:"type"`            // ??
}
