package order

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.Order, err error)
	Fetch(p r.Page) (res []*m.Order, err error)
	GetByID(id uint) (*m.Order, error)
	Update(id uint, v *m.Order) error
	Create(a *m.Order) error
	DeleteByID(id uint) error
	Delete(a *m.Order) error
}
