package model

import "github.com/jinzhu/gorm"

//ExEvent ...
type ExEvent struct {
	gorm.Model

	JadepoolID uint `gorm:"index;not null" json:"jadepoolID"`

	Log string `gorm:"type:json;not null" json:"log"` // store the whole json result into Log

	AssetID uint   `json:"assetID"`
	Hash    string `gorm:"index;type:varchar(128)" json:"hash"`
	Status  string `gorm:"type:varchar(32)" json:"status"` // PENDING, DONE, FAILED
}

//UpdateColumns ...
func (a *ExEvent) UpdateColumns(b *ExEvent) error {
	return GetDB().Model(ExEvent{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *ExEvent) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *ExEvent) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *ExEvent) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
