package address

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.Address, err error)
	Fetch(p r.Page) (res []*m.Address, err error)
	GetByID(id uint) (*m.Address, error)
	Update(id uint, v *m.Address) error
	Create(a *m.Address) error
	DeleteByID(id uint) error
	Delete(a *m.Address) error
}
