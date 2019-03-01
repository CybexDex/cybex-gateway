package jporder

import (
	m "coding.net/bobxuyang/cy-gateway-BN/models"
	r "coding.net/bobxuyang/cy-gateway-BN/repository"

	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.JPOrder, error)
	Fetch(p r.Page) ([]*m.JPOrder, error)
	FetchWith(o *m.JPOrder) ([]*m.JPOrder, error)
	GetByJPID(id uint) (*m.JPOrder, error)
	GetByID(id uint) (*m.JPOrder, error)
	DeleteByID(id uint) error
	Create(a *m.JPOrder) (err error)
	UpdateColumns(id uint, b *m.JPOrder) (err error)
	HoldingOne() (*m.JPOrder, error)
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
func (repo *Repo) FetchAll() ([]*m.JPOrder, error) {
	var res []*m.JPOrder
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.JPOrder, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.JPOrder) (res []*m.JPOrder, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.JPOrder, error) {
	a := m.JPOrder{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//GetByJPID ...
func (repo *Repo) GetByJPID(id uint) (*m.JPOrder, error) {
	a := m.JPOrder{}
	err := repo.DB.Where(&m.JPOrder{JadepoolOrderID: id}).First(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.JPOrder{}).Error
}

//Create ...
//for transaction base use
func (repo *Repo) Create(a *m.JPOrder) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		return err
	}

	return nil
}

//UpdateColumns ...
//for transaction base use
func (repo *Repo) UpdateColumns(id uint, b *m.JPOrder) error {
	return repo.DB.Model(m.JPOrder{}).Where("ID=?", id).UpdateColumns(b).Error
}

//HoldingOne ...
func (repo *Repo) HoldingOne() (*m.JPOrder, error) {
	var order m.JPOrder
	s := `update jp_orders 
	set status = 'HOLDING'
	where id = (
				select id 
				from jp_orders 
				where status = 'INIT' 
				order by id
				limit 1
			)
	returning *`
	err := repo.DB.Raw(s).Scan(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}
