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
func (repo *Repo) FetchAll() (res []*m.ExEvent, err error) {
	err = repo.DB.Find(&res).Error

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.ExEvent, error) {
	a := m.ExEvent{}
	err := repo.DB.First(&a, id).Error

	return &a, err
}

//Update ...
func (repo *Repo) Update(a *m.ExEvent) error {
	return nil
}

//Create ...
func (repo *Repo) Create(a *m.ExEvent) (err error) {
	return repo.DB.Create(&a).Error
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) (err error) {
	return repo.DB.Where("ID=?", id).Delete(&m.ExEvent{}).Error
}

//Delete ...
func (repo *Repo) Delete(a *m.ExEvent) (err error) {
	return repo.DB.Delete(&a).Error
}
