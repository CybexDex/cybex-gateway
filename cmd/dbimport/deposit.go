package main

import (
	"fmt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/cockroachdb/apd"
)

func main() {
	tOrder()
}

func tOrder() {
	jporderEntity := new(m.JPOrder)
	jporderEntity.From = "3QQDiUoKwNUVVnRY5Cyt5gKDhcocL7w5YP"
	jporderEntity.To = "1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1j"
	jporderEntity.Hash = "cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b"
	jporderEntity.UUHash = "BTC:cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b:1"
	jporderEntity.Index = 1
	jporderEntity.JadepoolOrderID = uint(408)
	jporderEntity.Status = "DONE"
	jporderEntity.Type = "DEPOSIT"
	jporderEntity.AssetID = 1
	jporderEntity.JadepoolID = 1
	amount, _, _ := apd.NewFromString("0.01000000")
	jporderEntity.Amount = amount
	// err := jporderEntity.Create()
	// if err != nil {
	// 	fmt.Println("jporderEntity.Create", err)
	// 	return
	// }

	orderEntity := new(m.Order)
	// orderEntity.JPHash = jporderEntity.Hash
	orderEntity.Status = "INIT"
	orderEntity.Type = jporderEntity.Type
	// orderEntity.JPUUHash = jporderEntity.UUHash
	orderEntity.AssetID = 1
	orderEntity.TotalAmount = amount
	// orderEntity.Amount = amount
	// fee, _, _ := apd.NewFromString("0")
	// orderEntity.Fee = fee
	orderEntity.AppID = 1
	err := orderEntity.Create()
	if err != nil {
		fmt.Println("orderEntity", err)
		return
	}
}
