package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Address ...
type Address struct {
	gorm.Model

	User       string `json:"user"`
	Asset      string `json:"asset"`
	BlockChain string `json:"blockChain"`
	Address    string `gorm:"index;type:varchar(128);not null" json:"address"`
}

// AddressLast ...
func AddressLast(user string, asset string) (address *Address, err error) {
	address = &Address{}
	if db == nil {
		return address, fmt.Errorf("no db init")
	}
	err = db.Where(&Address{
		Asset: asset,
		User:  user,
	}).Last(&address).Error
	if err != nil {
		return address, err
	}
	return address, nil
}

//AddrssCreate ...
func AddrssCreate(a *Address) error {
	err := db.Create(a).Error
	if err != nil {
		return err
	}
	return nil
}
