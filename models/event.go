package model

import "github.com/jinzhu/gorm"

//Event ...
type Event struct {
	gorm.Model

	AccountID uint `json:"accountID"`

	Route      string `gorm:"index;not null" json:"route"`
	Type       string `json:"type"` // login / change password / get address ...
	StatusCode int    `gorm:"not null" json:"statusCode"`
	UserAgent  string `json:"userAgent"`
	IPAddress  string `json:"ipAddress"`
	Input      string `json:"input"`
	Output     string `gorm:"type:json" json:"output"`
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
