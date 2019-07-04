package modeladmin

import (
	"github.com/jinzhu/gorm"
)

//Easy ...
type Easy struct {
	gorm.Model

	Key   string `gorm:"unique;index;type:varchar(128);default:null" json:"key"` // n to n, ???
	Value string `json:"value"`                                                  // people to event list
}

// Save ...
func (j *Easy) Save() error {
	return db.Save(j).Error
}

// EasyFristOrCreate ...
func EasyFristOrCreate(name string) (res *Easy, err error) {
	out := Easy{}
	err = db.Where(&Easy{
		Key: name,
	}).FirstOrCreate(&out).Error
	return &out, err
}

// EasyFind ...
func EasyFind(j *Easy) (res []*Easy, err error) {
	err = db.Where(j).Find(&res).Error
	return res, err
}
