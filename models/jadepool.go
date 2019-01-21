package model

import "github.com/jinzhu/gorm"

//Jadepool ...
type Jadepool struct {
	gorm.Model

	Apps     []App     `json:"apps"`     // 1 to n
	ExOrders []ExOrder `json:"exOrders"` // 1 to n
	ExEvents []ExEvent `json:"exEvents"` // 1 to n

	Name        string `gorm:"index;type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	TestNet     bool   `gorm:"default:false;not null" json:"testNet"`
	EccEnabled  bool   `gorm:"default:false;not null" json:"eccEnabled"`
	EccPubKey   string `gorm:"type:varchar(255)" json:"eccPubKey"`
	Host        string `gorm:"type:varchar(255);not null" json:"host"`
	Port        uint   `gorm:"not null" json:"port"`
	Version     string `gorm:"type:varchar(64)" json:"version"`
	Status      string `gorm:"type:varchar(32);not null" json:"status"` // INIT, NORMAL, ABNORMAL
	Type        string `gorm:"type:varchar(32)" json:"type"`            // DEFAULT, SAAS-VIP
}
