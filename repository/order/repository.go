package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.Order, err error)
	GetByID(id uint) (*m.Order, error)
	Update(a *m.Order) error
	Create(a *m.Order) error
	DeleteByID(id uint) error
	Delete(a *m.Order) error
}
