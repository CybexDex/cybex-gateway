package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

type Repository interface {
	FetchAll() (res []*m.Blockchain, err error)
	GetByID(id uint) (*m.Blockchain, error)
	Update(a *m.Blockchain) error
	Create(a *m.Blockchain) error
	DeleteByID(id uint) error
	Delete(a *m.Blockchain) error
}
