package bigasset

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.BigAsset, error)
	Fetch(p r.Page) ([]*m.BigAsset, error)
	FetchWith(o *m.BigAsset) ([]*m.BigAsset, error)
	GetByID(id uint) (*m.BigAsset, error)
	DeleteByID(id uint) error
	FetchWithOr(o1 *m.BigAsset, o2 *m.BigAsset) ([]*m.BigAsset, error)
	FindBig(query *m.BigAsset, big *apd.Decimal) (res []*m.BigAsset, err error)
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
func (repo *Repo) FetchAll() ([]*m.BigAsset, error) {
	var res []*m.BigAsset
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.BigAsset, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.BigAsset) (res []*m.BigAsset, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWithOr ...
func (repo *Repo) FetchWithOr(o1 *m.BigAsset, o2 *m.BigAsset) (res []*m.BigAsset, err error) {
	err = repo.DB.Where(o1).Or(o2).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FindBig ...
func (repo *Repo) FindBig(query *m.BigAsset, big *apd.Decimal) (res []*m.BigAsset, err error) {
	err = repo.DB.Where(query).Where("big_amount < ?", big).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.BigAsset, error) {
	a := m.BigAsset{}
	err := repo.DB.Preload("Asset").First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.BigAsset{}).Error
}
