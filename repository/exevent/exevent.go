package exorder

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/jinzhu/gorm"
)

//Repo ...
type Repo struct {
	DB *gorm.DB
}

//NewRepo ...
func NewRepo(db *gorm.DB) Repository {
	return &Repo{
		DB: db,
	}
}

//FetchAll ...
func (repo *Repo) FetchAll() (res []*m.ExEvent, err error) {
	err = repo.DB.Find(&res).Error
	if err != nil {
		utils.Errorf("ExEvent repo FetchAll error: ", err.Error())
		return nil, err
	}

	return res, nil
}

//GetByID ...
func (repo *Repo) GetByID(id uint) (a *m.ExEvent, err error) {
	err = repo.DB.First(&a, id).Error
	if err != nil {
		utils.Errorf("ExEvent repo GetByID error: ", err.Error())
		return nil, err
	}

	return a, nil
}

//Update ...
func (repo *Repo) Update(a *m.ExEvent) error {
	return nil
}

//Create ...
func (repo *Repo) Create(a *m.ExEvent) (err error) {
	err = repo.DB.Create(&a).Error
	if err != nil {
		utils.Errorf("ExEvent repo Create error: ", err.Error())
		return err
	}

	return nil
}

//DeleteByID ...
func (repo *Repo) DeleteByID(id uint) (err error) {
	err = repo.DB.Where("ID=?", id).Delete(&m.ExEvent{}).Error
	if err != nil {
		utils.Errorf("ExEvent repo DeleteByID error: ", err.Error())
		return err
	}

	return nil
}

//Delete ...
func (repo *Repo) Delete(a *m.ExEvent) (err error) {
	err = repo.DB.Delete(&a).Error
	if err != nil {
		utils.Errorf("ExEvent repo Delete error: ", err.Error())
		return err
	}

	return nil
}
