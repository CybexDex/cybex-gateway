package order

import (
	"database/sql"

	m "coding.net/bobxuyang/cy-gateway-BN/models"
	r "coding.net/bobxuyang/cy-gateway-BN/repository"
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
	UpdateAll(where *m.Order, update *m.Order) *gorm.DB
	HoldingOne() *m.Order
	QueryRecord(a *m.RecordsQuery) (out []*m.RecordsOut, err error)
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

// UpdateAll ...
func (repo *Repo) UpdateAll(where *m.Order, update *m.Order) *gorm.DB {
	return repo.DB.Model(&m.Order{}).Where(where).Update(update)
}

//HoldingOne ...
func (repo *Repo) HoldingOne() *m.Order {
	var order1 m.Order
	s := `update orders 
	set status = 'PROCESSING' 
	where id = (
				select id 
				from orders 
				where status = 'INIT' 
				order by id
				limit 1
			)
	returning *`
	repo.DB.Raw(s).Scan(&order1)
	return &order1
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

// QueryRecord ...
func (repo *Repo) QueryRecord(a *m.RecordsQuery) (resnew []*m.RecordsOut, err error) {
	res := []*m.Order{}
	err = repo.DB.Where(&m.Order{
		AppID: a.AppID,
		Type:  a.FundType,
	}).Preload("JPOrder").Preload("CybOrder").Preload("Asset").Preload("App").Offset(a.Offset).Limit(a.Size).Find(&res).Error
	if err != nil {
		return nil, err
	}
	// map
	for _, res1 := range res {
		var addr string
		if res1.Type == m.OrderTypeDeposit {
			addr = res1.JPOrder.From
		}
		if res1.Type == m.OrderTypeWithdraw {
			addr = res1.CybOrder.WithdrawAddr
		}
		resnew = append(resnew, &m.RecordsOut{
			Order:       res1,
			Asset:       res1.Asset.Name,
			CybexName:   res1.App.CybAccount,
			OutAddr:     addr,
			OutHash:     res1.JPOrder.Hash,
			CybHash:     res1.CybOrder.Hash,
			TotalAmount: res1.TotalAmount.Text('f'),
			Amount:      res1.Amount.Text('f'),
			Fee:         res1.Fee.Text('f'),
		})
	}
	return resnew, err
}
