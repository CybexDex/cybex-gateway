package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.Accounting, err error)
	GetByID(id uint) (*m.Accounting, error)
	Update(a *m.Accounting) error
	Create(a *m.Accounting) error
	DeleteByID(id uint) error
	Delete(a *m.Accounting) error
}
