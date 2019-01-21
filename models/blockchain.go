package model

import "github.com/jinzhu/gorm"

//Blockchain ...
type Blockchain struct {
	gorm.Model

	Assets []Asset `json:"assets"` // 1 to n

	Name         string `gorm:"unique;index;type:varchar(32);not null" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	Confirmation uint   `gorm:"default:20;not null" json:"confirmation"`
}
