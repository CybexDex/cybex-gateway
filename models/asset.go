package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Asset ...
type Asset struct {
	gorm.Model

	ExEvents  []ExEvent `json:"exEvents"`                            // 1 to n
	JPOrders  []JPOrder `json:"jpOrders"`                            // 1 to n
	Orders    []Order   `json:"orders"`                              // 1 to n
	Addresses []Address `gorm:"ForeignKey:AssetId" json:"addresses"` // 1 to n
	Balances  []Balance `gorm:"ForeignKey:AssetId" json:"balances"`  // 1 to n

	BlockchainID uint       `json:"blockchainID"` // n to 1
	Blockchain   Blockchain `gorm:"ForeignKey:BlockchainID" json:"blockchain"`

	Name           string       `gorm:"index;type:varchar(32);not null" json:"name"`
	Description    string       `gorm:"type:text" json:"description"`
	SmartContract  string       `gorm:"type:varchar(255)" json:"smartContract"`
	WithdrawSwitch bool         `json:"withdrawSwith"`
	DepositSwitch  bool         `json:"depositSwitch"`
	WithdrawFee    *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"withdrawFee"`
	DepositFee     *apd.Decimal `gorm:"type:numeric(30,10);default:0.0;not null" json:"depositFee"`
	LowWaterLevel  *apd.Decimal `gorm:"type:numeric(30,10)" json:"lowWaterLevel"`  // use for GATEWAY mode
	HighWaterLevel *apd.Decimal `gorm:"type:numeric(30,10)" json:"highWaterLevel"` // use for GATEWAY mode
	SweepTo        *apd.Decimal `gorm:"type:numeric(30,10)" json:"SweepTo"`        // use for GATEWAY mode
	Decimal        uint         `json:"decimal"`
}

//UpdateColumns ...
func (a *Asset) UpdateColumns(b *Asset) error {
	return GetDB().Model(Asset{}).Where("ID=?", a.ID).UpdateColumns(b).Error
}

//Create ...
func (a *Asset) Create() (err error) {
	return GetDB().Create(&a).Error
}

//Save ...
func (a *Asset) Save() (err error) {
	return GetDB().Save(&a).Error
}

//Delete ...
func (a *Asset) Delete() (err error) {
	return GetDB().Delete(&a).Error
}
