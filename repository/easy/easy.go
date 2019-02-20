package easy

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.Easy, error)
	Fetch(p r.Page) ([]*m.Easy, error)
	FetchWith(o *m.Easy) ([]*m.Easy, error)
	GetByID(id uint) (*m.Easy, error)
	DeleteByID(id uint) error
	FetchWithOr(o1 *m.Easy, o2 *m.Easy) ([]*m.Easy, error)
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
func (repo *Repo) FetchAll() ([]*m.Easy, error) {
	var res []*m.Easy
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.Easy, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.Easy) (res []*m.Easy, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWithOr ...
func (repo *Repo) FetchWithOr(o1 *m.Easy, o2 *m.Easy) (res []*m.Easy, err error) {
	err = repo.DB.Where(o1).Or(o2).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.Easy, error) {
	a := m.Easy{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.Easy{}).Error
}
