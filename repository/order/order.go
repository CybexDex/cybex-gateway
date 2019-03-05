package order

import (
	"database/sql"
	"fmt"

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
	QueryRecordAssets(a *m.RecordsQuery) (out []*RecordAssets, err error)
	FetchOrders(status string, fromnow string, offset int, limit int) (out []*m.Order, err error)
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

//FetchOrders ...
func (repo *Repo) FetchOrders(status string, fromnow string, offset int, limit int) (out []*m.Order, err error) {
	var s string
	if status != m.OrderStatusDone &&
		status != m.OrderStatusFailed &&
		status != m.OrderStatusInit &&
		status != m.OrderStatusPreInit &&
		status != m.OrderStatusProcessing &&
		status != m.OrderStatusTerminated &&
		status != m.OrderStatusWaiting {
		return nil, fmt.Errorf("status is invalid")
	}
	if len(fromnow) > 5 {
		return nil, fmt.Errorf("fromnow is invalid")
	}
	if fromnow == "" {
		s = fmt.Sprintf(`select * from orders where status = '%s'  order by id desc offset %d limit %d;`, status, offset, limit)
	} else {
		s = fmt.Sprintf(`select * from orders where status = '%s' and updated_at + interval '%s' < now()  order by id desc offset %d limit %d;`, status, fromnow, offset, limit)
	}
	err = repo.DB.Debug().Raw(s).Scan(&out).Error
	return out, err
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
	}).Preload("JPOrder").Preload("CybOrder").Preload("Asset").Preload("App").Order("id desc").Offset(a.Offset).Limit(a.Size).Find(&res).Error
	if err != nil {
		return nil, err
	}
	// map
	for _, res1 := range res {
		var addr string
		status := "PENDING"
		if res1.Type == m.OrderTypeDeposit {
			if res1.JPOrder != nil {
				addr = res1.JPOrder.From
			}
			if res1.CybOrder != nil {
				if res1.CybOrder.Status == m.CybOrderStatusDone {
					status = res1.CybOrder.Status
				}
			}
		}
		if res1.Type == m.OrderTypeWithdraw {
			if res1.CybOrder != nil {
				addr = res1.CybOrder.WithdrawAddr
			}
			if res1.JPOrder != nil {
				if res1.JPOrder.Status == m.JPOrderStatusDone {
					status = res1.JPOrder.Status
				}
			}
		}
		reout := &m.RecordsOut{
			Order:       res1,
			OutAddr:     addr,
			TotalAmount: res1.TotalAmount.Text('f'),
			Amount:      res1.Amount.Text('f'),
			Fee:         res1.Fee.Text('f'),
			Status:      status,
		}
		if res1.Asset != nil {
			reout.Asset = res1.Asset.Name
		}
		if res1.App != nil {
			reout.CybexName = res1.App.CybAccount
		}
		if res1.JPOrder != nil {
			reout.OutHash = res1.JPOrder.Hash
		}
		if res1.CybOrder != nil {
			reout.CybHash = res1.CybOrder.Hash
		}
		resnew = append(resnew, reout)
	}
	return resnew, err
}

// RecordAssets ...
type RecordAssets struct {
	Name  string `json:"asset"`
	Total int64  `json:"total"`
}

// QueryRecordAssets ...
func (repo *Repo) QueryRecordAssets(a *m.RecordsQuery) (out []*RecordAssets, err error) {
	s := fmt.Sprintf(`select assets.name,sum(1) as total from orders,assets where orders.asset_id = assets.id and  orders.app_id =%d group by assets.name;`, a.AppID)
	err = repo.DB.Raw(s).Scan(&out).Error
	return out, err
}
