package accounting

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.Accounting, err error)
	Fetch(p r.Page) (res []*m.Accounting, err error)
	GetByID(id uint) (*m.Accounting, error)
	Update(id uint, v *m.Accounting) error
	Create(a *m.Accounting) error
	DeleteByID(id uint) error
	Delete(a *m.Accounting) error
}
