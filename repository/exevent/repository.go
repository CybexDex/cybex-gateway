package exevent

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	r "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

//Repository ...
type Repository interface {
	FetchAll() (res []*m.ExEvent, err error)
	Fetch(p r.Page) (res []*m.ExEvent, err error)
	GetByID(id uint) (*m.ExEvent, error)
	Update(a *m.ExEvent) error
	Create(a *m.ExEvent) error
	DeleteByID(id uint) error
	Delete(a *m.ExEvent) error
}
