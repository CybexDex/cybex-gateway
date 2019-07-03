package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Task ...
type Task struct {
	gorm.Model

	Key    string `gorm:"index;type:varchar(128);default:null" json:"key"` // n to n, ???
	Value  string `json:"value"`                                           // people to event list
	Status string `json:"status"`
	Adds1  string `json:"adds1"`
	Adds2  string `json:"adds2"`
}

// Save ...
func (j *Task) Save() error {
	return db.Save(j).Error
}

// TaskFristOrCreate ...
func TaskFristOrCreate(name string) (res *Task, err error) {
	out := Task{}
	err = db.Where(&Task{
		Key: name,
	}).FirstOrCreate(&out).Error
	return &out, err
}

// WxSendTaskCreate ...
func WxSendTaskCreate(title string, msg string) (err error) {
	v := fmt.Sprintf("%s\n%s\n", title, msg)
	out := &Task{
		Key:    "wx",
		Value:  v,
		Status: "INIT",
	}
	return out.Save()
}

// HoldWxOne ...
func HoldWxOne() (*Task, error) {
	var order1 Task
	s := `update tasks 
	set status = 'PROCESSING' 
	where id = (
				select id 
				from tasks 
				where status = 'INIT' 
				and key = 'wx'
				order by id
				limit 1
			)
	returning *`
	err := db.Raw(s).Scan(&order1).Error
	return &order1, err
}

// TaskFind ...
func TaskFind(j *Task) (res []*Task, err error) {
	err = db.Where(j).Find(&res).Error
	return res, err
}
