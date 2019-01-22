package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.ExOrder, err error)
	GetByID(id uint) (*m.ExOrder, error)
	Update(a *m.ExOrder) error
	Create(a *m.ExOrder) error
	DeleteByID(id uint) error
	Delete(a *m.ExOrder) error
}
