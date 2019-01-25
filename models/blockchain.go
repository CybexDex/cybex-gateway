package model

import "github.com/jinzhu/gorm"

//Blockchain ...
type Blockchain struct {
	gorm.Model

	Assets []Asset `gorm:"ForeignKey:BlockchainId" json:"assets"` // 1 to n

	Name         string `gorm:"unique;index;type:varchar(32);not null" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	Confirmation uint   `gorm:"default:20;not null" json:"confirmation"`
}

//UpdateColumns ...
func (a *Blockchain) UpdateColumns(b *Blockchain) error {
	return GetDB().Model(Blockchain{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Blockchain) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Blockchain) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Blockchain) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
