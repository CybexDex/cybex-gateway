package model

import (
	"github.com/jinzhu/gorm"
)

//Easy ...
type Easy struct {
	gorm.Model

	Key   string `gorm:"unique;index;type:varchar(128);default:null" json:"key"` // n to n, ???
	Value string `json:"value"`                                                  // people to event list
}

//Save ...
func (a *Easy) Save() (err error) {
	return GetDB().Save(&a).Error
}
