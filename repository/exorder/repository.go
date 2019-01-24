package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() ([]*m.ExOrder, error)
	Fetch(p r.Page) ([]*m.ExOrder, error)
	FetchWith(o *m.ExOrder) ([]*m.ExOrder, error)
	GetByID(id uint) (*m.ExOrder, error)
	GetByJPID(id uint) (*m.ExOrder, error)
	Update(id uint, o *m.ExOrder) error
	Create(a *m.ExOrder) error
	DeleteByID(id uint) error
	Delete(a *m.ExOrder) error
}
