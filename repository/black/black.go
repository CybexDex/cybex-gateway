package black

import (
	m "coding.net/bobxuyang/cy-gateway-BN/models"
	r "coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.Black, error)
	Fetch(p r.Page) ([]*m.Black, error)
	FetchWith(o *m.Black) ([]*m.Black, error)
	GetByID(id uint) (*m.Black, error)
	DeleteByID(id uint) error
	FetchWithOr(o1 *m.Black, o2 *m.Black) ([]*m.Black, error)
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
func (repo *Repo) FetchAll() ([]*m.Black, error) {
	var res []*m.Black
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.Black, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.Black) (res []*m.Black, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWithOr ...
func (repo *Repo) FetchWithOr(o1 *m.Black, o2 *m.Black) (res []*m.Black, err error) {
	err = repo.DB.Where(o1).Or(o2).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.Black, error) {
	a := m.Black{}
	err := repo.DB.Preload("Asset").First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.Black{}).Error
}
