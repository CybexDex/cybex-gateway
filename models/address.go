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
