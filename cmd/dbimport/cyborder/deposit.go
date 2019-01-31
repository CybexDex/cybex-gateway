package cyborder

import (
	"fmt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/cockroachdb/apd"
)

func ToCYBOrders() {
	// tOrder()
	for i := 1; i <= 100; i++ {
		ToCYBOrder(i)
	}
}

func ToCYBOrder(i int) {
	var db = m.GetDB()

	// App
	totalamount, _, _ := apd.NewFromString("0.01000000")
	fee, _, _ := apd.NewFromString("0.00200000")
	s := fmt.Sprintf("0.001%d", i)
	amount, _, _ := apd.NewFromString(s)
	app := m.CybOrder{
		From: "yangyu123",
		To:   "yangyu1",
		// FeeAssetID uint

		TotalAmount: totalamount, // totalAmount = amount + fee
		Amount:      amount,      // `gorm:"type:numeric(30,10);not null" json:"amount"`
		Fee:         fee,         // `gorm:"type:numeric(30,10);not null" json:"fee"` // fee in Asset
		AssetID:     1,
		Status:      "INIT",    // INIT, HOLDING, PENDING, DONE, FAILED
		Type:        "DEPOSIT", // string `gorm:"type:varchar(32);not null" json:"type"`                 // DEPOSIT, WITHDRAW, RECHARGE, SWEEP, FEESETTLE
		Settled:     false,     //`gorm:"not null;default:false" json:"settled"`
	}
	err := db.Create(&app).Error
	fmt.Println(err)
}
