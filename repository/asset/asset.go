package asset

import (
	"fmt"

	m "coding.net/bobxuyang/cy-gateway-BN/models"
	r "coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.Asset, error)
	Fetch(p r.Page) ([]*m.Asset, error)
	FetchWith(o *m.Asset) ([]*m.Asset, error)
	GetByName(name string) (*m.Asset, error)
	GetByID(id uint) (*m.Asset, error)
	DeleteByID(id uint) error
	RawAmountToReal(amount int64, precision int) *apd.Decimal
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
func (repo *Repo) FetchAll() ([]*m.Asset, error) {
	var res []*m.Asset
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.Asset, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.Asset) (res []*m.Asset, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.Asset, error) {
	a := m.Asset{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByName ...
func (repo *Repo) GetByName(name string) (*m.Asset, error) {
	a := m.Asset{}
	err := repo.DB.Where("name=?", name).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.Asset{}).Error
}

// RawAmountToReal ...
func (repo *Repo) RawAmountToReal(amount int64, precision int) *apd.Decimal {
	// precics
	amountStr, _, _ := apd.NewFromString(fmt.Sprintf("%d", uint64(amount)))
	p, _, _ := apd.NewFromString(fmt.Sprintf("%d", precision))
	p1, _, _ := apd.NewFromString("10")
	out, _, _ := apd.NewFromString("0")
	c := apd.BaseContext
	c.Precision = 10
	c.Pow(p, p1, p)
	c.Quo(out, amountStr, p)
	return out
}
