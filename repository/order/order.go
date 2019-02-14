package order

import (
	"database/sql"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.Order, error)
	Fetch(p r.Page) ([]*m.Order, error)
	FetchWith(o *m.Order) ([]*m.Order, error)
	GetByID(id uint) (*m.Order, error)
	DeleteByID(id uint) error
	Create(a *m.Order) (err error)
	Rows(o *m.Order) (*sql.Rows, error)
	MDB() *gorm.DB
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
func (repo *Repo) FetchAll() ([]*m.Order, error) {
	var res []*m.Order
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.Order, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.Order) (res []*m.Order, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.Order, error) {
	a := m.Order{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.Order{}).Error
}

// MDB ...
func (repo *Repo) MDB() *gorm.DB {
	return repo.DB.Model(&m.Order{})
}

//Rows
func (repo *Repo) Rows(o *m.Order) (*sql.Rows, error) {
	rows, err := repo.DB.Where(o).Rows()
	defer rows.Close()
	return rows, err
}

//Create ...
//for transaction base use
func (repo *Repo) Create(a *m.Order) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		return err
	}

	return nil
}
