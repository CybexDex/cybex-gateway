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

//FetchAll ...
func (blockchain *Blockchain) FetchAll() ([]*Blockchain, error) {
	var res []*Blockchain
	err := GetDB().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (blockchain *Blockchain) Fetch(p Page) (res []*Blockchain, err error) {
	err = GetDB().Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (blockchain *Blockchain) GetByID(id uint) (*Blockchain, error) {
	a := Blockchain{}
	err := GetDB().First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//Update ...
func (blockchain *Blockchain) Update(id uint, v *Blockchain) error {
	return GetDB().Model(Blockchain{}).Where("ID=?", id).UpdateColumns(v).Error
}

//Create ...
func (blockchain *Blockchain) Create(a *Blockchain) (err error) {
	return GetDB().Create(&a).Error
}

//DeleteByID ...
func (blockchain *Blockchain) DeleteByID(id uint) (err error) {
	return GetDB().Where("ID=?", id).Delete(&Blockchain{}).Error
}

//Delete ...
func (blockchain *Blockchain) Delete(a *Blockchain) (err error) {
	return GetDB().Delete(&a).Error
}
