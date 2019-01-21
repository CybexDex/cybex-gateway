package model

import "github.com/jinzhu/gorm"

//ExEvent ...
type ExEvent struct {
	gorm.Model

	JadepoolID uint `gorm:"index;not null" json:"jadepoolID"`

	Log string `gorm:"type:text;not null" json:"log"` // store the whole json result into Log

	AssetID uint   `json:"assetID"`
	Hash    string `gorm:"index;type:varchar(128)" json:"hash"`
	Status  string `gorm:"type:varchar(32)" json:"status"` // PENDING, DONE, FAILED
}
