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
	Blockchain string `gorm:"type:varchar(32);not null" json:"blockchain"`
}

// BlackWith ...
func BlackWith(user string, blockchain string, addresses ...string) (bool, []*Black, error) {
	var bs []*Black
	err := db.Where(&Black{
		Address:    user,
		Blockchain: "CYB",
	}).Find(&bs).Error
	if err != nil {
		return true, nil, err
	}
	if len(bs) > 0 {
		return true, bs, nil
	}
	err = db.Where("blockchain = ? and address in (?)", blockchain, addresses).Find(&bs).Error
	if err != nil {
		return true, nil, err
	}
	if len(bs) > 0 {
		return true, bs, nil
	}
	return false, nil, nil
}
