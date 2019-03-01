package cyborder

import (
	"fmt"
	"math/rand"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/cockroachdb/apd"
	"github.com/lib/pq"
)

func ToCYBOrders() {
	// tOrder()
	for i := 1; i <= 10; i++ {
		// ToCYBOrder(i)
	}
	ToBigAsset()
	ToBlack()
	ToOrder()
}
func ToOrder() {
	// var db = m.GetDB()
	amount, _, _ := apd.NewFromString("0.01")
	orderEntity := new(m.Order)
	orderEntity.JPOrderID = 1
	orderEntity.CybOrderID = 1
	orderEntity.FailedJPOrders = *new(pq.Int64Array)
	orderEntity.FailedJPOrders = append(orderEntity.FailedJPOrders, 1)
	orderEntity.FailedJPOrders = append(orderEntity.FailedJPOrders, 1)
	orderEntity.FailedCybOrders = *new(pq.Int64Array)
	orderEntity.FailedCybOrders = append(orderEntity.FailedCybOrders, 1)
	orderEntity.FailedCybOrders = append(orderEntity.FailedCybOrders, 1)
	orderEntity.JPHash = "cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b1"
	orderEntity.JPUUHash = "BTC:cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b1:1"
	orderEntity.CybHash = "cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b"
	orderEntity.CybUUHash = "BTC:cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b:1"
	orderEntity.Status = "INIT"
	orderEntity.Type = "DEPOSIT"
	orderEntity.AssetID = 1
	orderEntity.TotalAmount = amount
	orderEntity.Amount = amount
	fee, _, _ := apd.NewFromString("0")
	orderEntity.Fee = fee
	orderEntity.AppID = 1
	err := orderEntity.Create()
	fmt.Println(err)
}
func ToBlack() {
	var db = m.GetDB()
	black := m.Black{
		Address:    "yangyu1",
		Blockchain: "CYB",
	}
	err := db.Create(&black).Error
	fmt.Println(err)
}
func ToBigAsset() {
	var db = m.GetDB()
	amount, _, _ := apd.NewFromString("10")
	app := m.BigAsset{
		AssetID:   1,
		BigAmount: amount,
		Type:      "DEPOSIT", // string `gorm:"type:varchar(32);not null" json:"type"`                 // DEPOSIT, WITHDRAW, RECHARGE, SWEEP, FEESETTLE
		Level:     1,         //`gorm:"not null;default:false" json:"settled"`
	}
	err := db.Create(&app).Error
	fmt.Println(err)
}
func ToCYBOrder(i int) {
	var db = m.GetDB()

	// App
	totalamount, _, _ := apd.NewFromString("0.01000000")
	fee, _, _ := apd.NewFromString("0.00200000")
	s := fmt.Sprintf("0.001%d", i)
	amount, _, _ := apd.NewFromString(s)
	tos := []string{"yangyu1", "yangyu2"}

	j := rand.Intn(len(tos))
	to := tos[j]
	app := m.CybOrder{
		From: "yangyu123",
		To:   to,
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
