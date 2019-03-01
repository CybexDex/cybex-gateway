package cyborder

import (
	m "coding.net/bobxuyang/cy-gateway-BN/models"
	r "coding.net/bobxuyang/cy-gateway-BN/repository"
	"github.com/jinzhu/gorm"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.CybOrder, error)
	Fetch(p r.Page) ([]*m.CybOrder, error)
	FetchWith(o *m.CybOrder) ([]*m.CybOrder, error)
	GetByID(id uint) (*m.CybOrder, error)
	DeleteByID(id uint) error
	Create(a *m.CybOrder) (err error)
	MDB() *gorm.DB
	HoldingOne() *m.CybOrder
	UpdateAll(where *m.CybOrder, update *m.CybOrder) *gorm.DB
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

// MDB ...
func (repo *Repo) MDB() *gorm.DB {
	return repo.DB.Model(&m.CybOrder{})
}

// UpdateAll ...
func (repo *Repo) UpdateAll(where *m.CybOrder, update *m.CybOrder) *gorm.DB {
	return repo.DB.Model(&m.CybOrder{}).Where(where).Update(update)
}

//HoldingOne ...
func (repo *Repo) HoldingOne() *m.CybOrder {
	var order1 m.CybOrder
	s := `update cyb_orders 
	set status = 'HOLDING' 
	where id = (
				select id 
				from cyb_orders 
				where status = 'INIT' 
				order by id
				limit 1
			)
	returning *`
	repo.DB.Raw(s).Scan(&order1)
	return &order1
}

//FetchAll ...
func (repo *Repo) FetchAll() ([]*m.CybOrder, error) {
	var res []*m.CybOrder
	err := repo.DB.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//Fetch ...
func (repo *Repo) Fetch(p r.Page) (res []*m.CybOrder, err error) {
	err = repo.DB.Order(p.OrderBy + " " + p.Sort).Offset(p.Offset).Find(&res).Limit(p.Amount).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//FetchWith ...
func (repo *Repo) FetchWith(o *m.CybOrder) (res []*m.CybOrder, err error) {
	err = repo.DB.Where(o).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, err
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (*m.CybOrder, error) {
	a := m.CybOrder{}
	err := repo.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}

	return &a, err
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) error {
	return repo.DB.Where("ID=?", id).Delete(&m.CybOrder{}).Error
}

//Create ...
//for transaction base use
func (repo *Repo) Create(a *m.CybOrder) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		return err
	}

	return nil
}
