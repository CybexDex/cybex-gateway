package app

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.App, error)
	Fetch(p r.Page) ([]*m.App, error)
	FetchWith(o *m.App) ([]*m.App, error)
	GetByName(name string) (*m.App, error)
	GetByID(id uint) (*m.App, error)
	DeleteByID(id uint) error
}

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
func (repo *Repo) FetchAll() ([]*m.App, error) {
	var res []*m.App
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.App, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.App) (res []*m.App, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.App, error) {
	a := m.App{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByName ...
func (repo *Repo) GetByName(name string) (*m.App, error) {
	a := m.App{}
	err := repo.DB.Where("name=?", name).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.App{}).Error
}
