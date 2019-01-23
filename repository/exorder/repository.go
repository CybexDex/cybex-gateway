package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.ExOrder, err error)
	Fetch(p r.Page) (res []*m.ExOrder, err error)
	GetByID(id uint) (*m.ExOrder, error)
	GetByJPID(id uint) (*m.ExOrder, error)
	Update(id uint, v *m.ExOrder) error
	Create(a *m.ExOrder) error
	DeleteByID(id uint) error
	Delete(a *m.ExOrder) error
}
