package modeladmin

import (
	"github.com/jinzhu/gorm"
)

//OrderLog ...
type OrderLog struct {
	gorm.Model
	OrderID uint   `gorm:"index" json:"orderID"`
	Event   string `json:"event"`
	Message string `json:"message"`
	Offset  int    `gorm:"-" json:"offset"`
	Limit   int    `gorm:"-" json:"limit"`
}

// ShowLogs ...
func ShowLogs(j *OrderLog) (res []*OrderLog, total int, err error) {
	err = db.Debug().Where(j).Order("id").Offset(j.Offset).Limit(j.Limit).Find(&res).Error
	if err != nil {
		return res, total, err
	}
	var x []*OrderLog
	err = db.Debug().Where(j).Find(&x).Count(&total).Error
	return res, total, err
}
