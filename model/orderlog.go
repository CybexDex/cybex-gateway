package model

import (
	"github.com/jinzhu/gorm"
)

//OrderLog ...
type OrderLog struct {
	gorm.Model
	OrderID uint   `json:"orderID"`
	Event   string `json:"event"`
	Message string `json:"message"`
}
