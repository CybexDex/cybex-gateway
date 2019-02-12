package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//BigAsset ...
type BigAsset struct {
	gorm.Model
	AssetID   uint         `gorm:"index:id_type_level" json:"assetID"`
	Type      string       `gorm:"type:varchar(32);not null;index:id_type_level" json:"type"` // DEPOSIT, WITHDRAW
	BigAmount *apd.Decimal `gorm:"type:numeric(30,10)" json:"bigAmount"`
	Level     uint         `gorm:"index:id_type_level" json:"level"` // 1,2,3 higher is bigger
}
