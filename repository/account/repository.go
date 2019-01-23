package account

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.Account, err error)
	Fetch(p r.Page) (res []*m.Account, err error)
	GetByName(name string) (*m.Account, error)
	GetByID(id uint) (*m.Account, error)
	Update(id uint, v *m.Account) error
	Create(a *m.Account) error
	DeleteByID(id uint) error
	Delete(a *m.Account) error
}
