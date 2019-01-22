package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.ExEvent, err error)
	GetByID(id uint) (*m.ExEvent, error)
	Update(a *m.ExEvent) error
	Create(a *m.ExEvent) error
	DeleteByID(id uint) error
	Delete(a *m.ExEvent) error
}
