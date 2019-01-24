package exevent

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.ExEvent, error)
	Fetch(p r.Page) ([]*m.ExEvent, error)
	FetchWith(o *m.ExEvent) ([]*m.ExEvent, error)
	Create(a *m.ExEvent) (err error)
	GetByName(name string) (*m.ExEvent, error)
	GetByID(id uint) (*m.ExEvent, error)
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
func (repo *Repo) FetchAll() ([]*m.ExEvent, error) {
	var res []*m.ExEvent
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.ExEvent, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.ExEvent) (res []*m.ExEvent, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.ExEvent, error) {
	a := m.ExEvent{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByName ...
func (repo *Repo) GetByName(name string) (*m.ExEvent, error) {
	a := m.ExEvent{}
	err := repo.DB.Where("name=?", name).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.ExEvent{}).Error
}

//Create ...
//for transaction base use
func (repo *Repo) Create(a *m.ExEvent) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		return err
	}

	return nil
}
