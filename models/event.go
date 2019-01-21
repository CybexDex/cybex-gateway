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
