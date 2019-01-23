package model

import (
	"github.com/jinzhu/gorm"
)

//Account ...
type Account struct {
	gorm.Model

	Apps   []App   `gorm:"many2many:account_apps;" json:"apps"` // n to n, ???
	Events []Event `json:"events"`                              // people to event list

	CompanyID uint `json:"companyID"`

	Name            string `gorm:"unique;index;type:varchar(64);not null" json:"name"`
	LastName        string `gorm:"type:varchar(32)" json:"lastName"`
	FirstName       string `gorm:"type:varchar(32)" json:"firstName"`
	Email           string `gorm:"unique;index;type:varchar(64);not null" json:"email"`
	EmailVerified   bool   `gorm:"default:false" json:"emailVerified"`
	EmailEnable     bool   `gorm:"default:false" json:"emailEnabled"`
	Phone           string `gorm:"type:varchar(32)" json:"phone"`
	PhoneVerified   bool   `gorm:"default:false" json:"phoneVerified"`
	PhoneEnabled    bool   `gorm:"default:false" json:"phoneEnabled"`
	AuthKey         string `gorm:"type:varchar(64)" json:"authKey"`
	AuthKeyVerified bool   `gorm:"default:false" json:"authKeyVerified"`
	AuthKeyEnabled  bool   `gorm:"default:false" json:"authKeyEnabled"`
	PasswordHash    string `gorm:"type:varchar(512);not null" json:"passwordHash"`
	Password        string `gorm:"-" json:"password"`
	Status          string `gorm:"type:varchar(32);default:'INIT'" json:"status"`    // INIT, NORMAL, ABNORMAL
	Type            string `gorm:"type:varchar(32);default:'SAAS_USER'" json:"type"` // ADMIN, SUPER_ADMIN, SAAS_USER
	Disable         bool   `gorm:"default:false" json:"disable"`
}

//FetchAll ...
func (acc *Account) FetchAll() ([]*Account, error) {
	var res []*Account
	err := GetDB().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (acc *Account) Fetch(p Page) (res []*Account, err error) {
	err = GetDB().Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (acc *Account) GetByID(id uint) (*Account, error) {
	a := Account{}
	err := GetDB().First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByName ...
func (acc *Account) GetByName(name string) (*Account, error) {
	a := Account{}
	err := GetDB().Where("name=?", name).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//Update ...
func (acc *Account) Update(id uint, v *Account) error {
	return GetDB().Model(Account{}).Where("ID=?", id).UpdateColumns(v).Error
}

//Create ...
func (acc *Account) Create(a *Account) (err error) {
	return GetDB().Create(&a).Error
}

//DeleteByID ...
func (acc *Account) DeleteByID(id uint) (err error) {
	return GetDB().Where("ID=?", id).Delete(&Account{}).Error
}

//Delete ...
func (acc *Account) Delete(a *Account) (err error) {
	return GetDB().Delete(&a).Error
}
