package model

import (
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//Asset ...
type Asset struct {
	gorm.Model

	ExEvents    []ExEvent    `json:"exEvents"`   // 1 to n
	ExOrders    []ExOrder    `json:"exOrders"`   // 1 to n
	Orders      []Order      `json:"orders"`     // 1 to n
	Addresses   []Address    `json:"addresses"`  // 1 to n
	Accountings []Accounting `json:"accounting"` // 1 to n
	Balances    []Balance    `json:"balances"`   // 1 to n

	BlockchainID uint `json:"blockchainID"` // n to 1

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
