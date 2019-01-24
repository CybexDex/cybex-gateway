package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"

	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.ExOrder, error)
	Fetch(p r.Page) ([]*m.ExOrder, error)
	FetchWith(o *m.ExOrder) ([]*m.ExOrder, error)
	GetByJPID(id uint) (*m.ExOrder, error)
	GetByName(name string) (*m.ExOrder, error)
	GetByID(id uint) (*m.ExOrder, error)
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

//FetchWith ...
func (repo *Repo) FetchWith(o *m.ExOrder) (res []*m.ExOrder, err error) {
	err = repo.DB.Where(o).Find(&res).Error
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

//GetByName ...
func (repo *Repo) GetByName(name string) (*m.ExOrder, error) {
	a := m.ExOrder{}
	err := repo.DB.Where("name=?", name).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.ExOrder{}).Error
}
