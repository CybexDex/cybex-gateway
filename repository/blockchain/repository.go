package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.Blockchain, err error)
	Fetch(p r.Page) (res []*m.Blockchain, err error)
	GetByID(id uint) (*m.Blockchain, error)
	Update(a *m.Blockchain) error
	Create(a *m.Blockchain) error
	DeleteByID(id uint) error
	Delete(a *m.Blockchain) error
}
