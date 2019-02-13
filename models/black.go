package model

import (
	"github.com/jinzhu/gorm"
)

//dbimport
//base

//Black ...
type Black struct {
	gorm.Model
	Address    string `gorm:"type:varchar(128)" json:"address"`
	Blockchain string `gorm:"unique;index;type:varchar(32);not null" json:"blockchain"`
}
