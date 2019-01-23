package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"

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
func (repo *Repo) FetchAll() ([]*m.ExOrder, error) {
	var res []*m.ExOrder
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.ExOrder, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.ExOrder, error) {
	a := m.ExOrder{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByJPID ...
func (repo *Repo) GetByJPID(id uint) (*m.ExOrder, error) {
	a := m.ExOrder{}
	err := repo.DB.Where(&m.ExOrder{JadepoolOrderID: id}).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//Update ...
func (repo *Repo) Update(id uint, v *m.ExOrder) error {
	return repo.DB.Model(m.ExOrder{}).Where("ID=?", id).UpdateColumns(v).Error
}

//Create ...
func (repo *Repo) Create(a *m.ExOrder) (err error) {
	return repo.DB.Create(&a).Error
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) (err error) {
	return repo.DB.Where("ID=?", id).Delete(&m.ExOrder{}).Error
}

//Delete ...
func (repo *Repo) Delete(a *m.ExOrder) (err error) {
	return repo.DB.Delete(&a).Error
}
