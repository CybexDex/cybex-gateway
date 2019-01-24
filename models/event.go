package model

import "github.com/jinzhu/gorm"

//Event ...
type Event struct {
	gorm.Model

	AccountID uint `json:"accountID"`

	Type      string `json:"type"` // login / change password / get address ...
	Input     string `json:"input"`
	Output    string `json:"output"`
	Result    int    `json:"result"` // same as http result code
	ResultStr string `gorm:"type:text"`
}

//UpdateColumns ...
func (a *Event) UpdateColumns(b *Event) error {
	return GetDB().Model(Event{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Event) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Event) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Event) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
