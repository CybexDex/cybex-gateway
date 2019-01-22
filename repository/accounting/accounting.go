package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/jinzhu/gorm"
)

//Repo ...
type Repo struct {
	DB *gorm.DB
}

//NewRepo ...
func NewRepo(db *gorm.DB) Repository {
	return &Repo{
		DB: db,
	}
}

//FetchAll ...
func (repo *Repo) FetchAll() (res []*m.Accounting, err error) {
	err = repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.Accounting, error) {
	a := m.Accounting{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, nil
}

//Update ...
func (repo *Repo) Update(a *m.Accounting) error {
	return nil
}

//Create ...
func (repo *Repo) Create(a *m.Accounting) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		return err
	}

	return nil
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) (err error) {
	err = repo.DB.Where("ID=?", id).Delete(&m.Accounting{}).Error
	if err != nil {
		return err
	}

	return nil
}

//Delete ...
func (repo *Repo) Delete(a *m.Accounting) (err error) {
	err = repo.DB.Delete(&a).Error
	if err != nil {
		return err
	}

	return nil
}
