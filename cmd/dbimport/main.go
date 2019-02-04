package main

import (
	"fmt"

	"git.coding.net/bobxuyang/cy-gateway-BN/cmd/dbimport/cyborder"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/cockroachdb/apd"
	"github.com/lib/pq"
)

func main() {
	tBlockchain()
	tAsset()
	tJadepool()
	tCompany()
	tAccount()
	tApp()
	tAddress()
	tExOrder()
	tQueryPreload()
	tQueryPreload2()
	tOrder1()
	tBalance()
	cyborder.ToCYBOrders()
}

func tBlockchain() {
	var db = m.GetDB()

	// blockchain
	bc := m.Blockchain{Name: "Ethereum", Description: "Ethereum", Confirmation: 30}
	err := db.Create(&bc).Error
	fmt.Println(err)

	bc = m.Blockchain{Name: "Bitcoin", Description: "Bitcoin", Confirmation: 6}
	err = db.Create(&bc).Error
	fmt.Println(err)
}

func tAsset() {
	var db = m.GetDB()

	// assets
	wf, _, _ := new(apd.Decimal).SetString("0.05")
	df, _, _ := new(apd.Decimal).SetString("0.05")
	lwl, _, _ := new(apd.Decimal).SetString("100.0")
	hwl, _, _ := new(apd.Decimal).SetString("10000.0")
	st, _, _ := new(apd.Decimal).SetString("1000.0")
	asset := m.Asset{
		Name:           "ETH",
		CybName:        "CYB",
		CybID:          "1.3.0",
		Description:    "ETH",
		BlockchainID:   1,
		WithdrawSwitch: true,
		DepositSwitch:  true,
		WithdrawFee:    wf,
		DepositFee:     df,
		LowWaterLevel:  lwl,
		HighWaterLevel: hwl,
		SweepTo:        st,
		Decimal:        18,
	}

	err := db.Create(&asset).Error
	fmt.Println(err)
}

func tJadepool() {
	var db = m.GetDB()

	// jadepool
	jp := m.Jadepool{
		Name:        "Jadepool-001",
		Description: "Jadepool-001",
		TestNet:     true,
		EccEnabled:  false,
		Host:        "192.168.0.1",
		Port:        6688,
		Version:     "1.0.0",
		Status:      "NORMAL",
		Type:        "DEFAULT",
	}
	err := db.Create(&jp).Error
	fmt.Println(err)
}

func tCompany() {
	var db = m.GetDB()

	// company
	company := m.Company{
		Name:        "nbltrust",
		Description: "nbltrust",
		Status:      "NORMAL",
		Type:        "VIP",
		CompanyAddress: m.GeoAddress{
			Address:   "Soho 3Q 2 Floor",
			Zipcode:   "200001",
			LastName:  "Yang",
			FirstName: "Yu",
			Phone:     "13813812345",
		},
		ContactAddress: m.GeoAddress{
			Address:   "Soho 3Q 2 Floor",
			Zipcode:   "200005",
			LastName:  "Xu",
			FirstName: "Yang",
			Phone:     "13912300000",
		},
	}
	err := db.Create(&company).Error
	fmt.Println(err)
}

func tAccount() {
	var db = m.GetDB()

	// Account
	account := m.Account{
		CompanyID: 1,
		// Name:      &sql.NullString{String: "xuyang", Valid: true},
		Name:            "xuyang",
		LastName:        "XU",
		FirstName:       "YANG",
		Email:           "xuyang@nbltrust.com",
		EmailEnable:     false,
		EmailVerified:   false,
		Phone:           "13912300000",
		PhoneEnabled:    false,
		PhoneVerified:   false,
		AuthKey:         "139123000001391230000013912300000",
		AuthKeyEnabled:  false,
		AuthKeyVerified: false,
		PasswordHash:    "04a24e8195382cbfe6c81dda873d2be49b13c1bd09b01f0bfeeba952de3c59cd",
		Disable:         false,
	}
	err := db.Create(&account).Error
	fmt.Println(err)
}

func tApp() {
	var db = m.GetDB()

	// App
	app := m.App{
		CompanyID:   1,
		JadepoolID:  1,
		Name:        "some awesome appplication",
		CybAccount:  "yangyu111",
		Description: "some awesome appplication",
		URL:         "https://www.app.dcom",
		Status:      "NORMAL",
		Type:        "professional",
	}
	err := db.Create(&app).Error
	fmt.Println(err)
}

func tAddress() {
	var db = m.GetDB()

	address := m.Address{
		AppID:   1,
		AssetID: 1,
		Address: "0xa48d73341885e6bce0252cb05b31a8a00720cdb2",
		Status:  "NORMAL",
	}
	err := db.Create(&address).Error
	fmt.Println(err)
}

func tExOrder() {
	var db = m.GetDB()

	a, _, _ := new(apd.Decimal).SetString("100.0")

	jporder := m.JPOrder{
		AssetID:         1,
		JadepoolID:      1,
		JadepoolOrderID: 100,
		From:            "0x5b19816a9be1aaea3664f583f9e5fd76188d1402",
		To:              "0xb45a9e7878e74117509653538d7bb7f4122352d2",
		Amount:          a,
		Index:           0,
		Hash:            "0xd949dd10db2c5a5c45b0f4b2899851783546ffe4d71f848d9b1505933d01cd37",
		UUHash:          "ETHEREUM:0xd949dd10db2c5a5c45b0f4b2899851783546ffe4d71f848d9b1505933d01cd37",
		Status:          "PENDING",
		Type:            "DEPOSIT",
	}
	err := db.Create(&jporder).Error
	if err != nil {
		fmt.Println("jporderEntity", err)
		return
	}
}

func tQueryPreload() {
	var db = m.GetDB()

	var com m.Company
	db.Where("id=?", 1).Preload("Accounts").Preload("Apps").Preload("CompanyAddress").Preload("ContactAddress").First(&com)
	fmt.Println(len(com.Accounts), com.Accounts)
	fmt.Println(len(com.Apps), com.Apps)
	fmt.Println(com.CompanyAddress)
	fmt.Println(com.ContactAddress)
	fmt.Println(com)
}

func tQueryPreload2() {
	var db = m.GetDB()

	var com m.Company
	db.Where("id=?", 1).Preload("Accounts").Preload("Apps").First(&com)
	fmt.Println(len(com.Accounts), com.Accounts)
	fmt.Println(len(com.Apps), com.Apps)
	fmt.Println(com)
}

func tOrder1() {
	var db = m.GetDB()

	zero := apd.New(0, 0)

	jporderEntity := new(m.JPOrder)
	jporderEntity.From = "3QQDiUoKwNUVVnRY5Cyt5gKDhcocL7w5YP"
	jporderEntity.To = "1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1j"
	jporderEntity.Hash = "cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b"
	jporderEntity.UUHash = "BTC:cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b:1"
	jporderEntity.Fee = zero
	jporderEntity.TotalAmount = zero
	jporderEntity.Amount = zero

	jporderEntity.Index = 1
	jporderEntity.JadepoolOrderID = uint(408)
	jporderEntity.Status = "DONE"
	jporderEntity.Type = "DEPOSIT"
	jporderEntity.AssetID = 1
	jporderEntity.AppID = 1
	jporderEntity.JadepoolID = 1
	amount, _, _ := apd.NewFromString("0.01000000")
	jporderEntity.Amount = amount
	err := db.Save(jporderEntity).Error
	if err != nil {
		fmt.Println("jporderEntity.Create", err)
		return
	}

	cyborderEntity := new(m.CybOrder)
	cyborderEntity.AssetID = 1
	cyborderEntity.AppID = 1
	cyborderEntity.From = "xuyang"
	cyborderEntity.To = "yangyu"
	amount, _, _ = apd.NewFromString("100")
	cyborderEntity.Amount = amount
	cyborderEntity.Hash = "400000:3"
	cyborderEntity.UUHash = "cyb:400000:3"
	cyborderEntity.Fee = zero
	cyborderEntity.Amount = zero
	cyborderEntity.TotalAmount = zero
	cyborderEntity.Status = "PENDING"
	cyborderEntity.Type = "DEPOSIT"
	err = cyborderEntity.Create()
	if err != nil {
		fmt.Println("cyborderEntity", err)
		return
	}

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
	orderEntity.Fee = zero
	orderEntity.Amount = zero
	orderEntity.TotalAmount = zero
	orderEntity.JPHash = jporderEntity.Hash
	orderEntity.Status = "INIT"
	orderEntity.Type = jporderEntity.Type
	orderEntity.JPUUHash = jporderEntity.UUHash
	orderEntity.AssetID = 1
	orderEntity.TotalAmount = amount
	orderEntity.Amount = amount
	fee, _, _ := apd.NewFromString("0")
	orderEntity.Fee = fee
	orderEntity.AppID = 1
	err = orderEntity.Create()
	if err != nil {
		fmt.Println("orderEntity...", err)
		return
	}
}

func tBalance() {
	bal := new(m.Balance)
	bal.AppID = 1
	bal.AssetID = 1
	err := bal.Save()

	if err != nil {
		fmt.Println("balance", err)
		return
	}
}
