package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.Address, err error)
	GetByID(id uint) (*m.Address, error)
	Update(a *m.Address) error
	Create(a *m.Address) error
	DeleteByID(id uint) error
	Delete(a *m.Address) error
}
