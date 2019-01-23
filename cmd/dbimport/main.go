package main

import (
	"fmt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/cockroachdb/apd"
)

func main() {
	tBlockchain()
	tAsset()
	tJadepool()
	tCompany()
	tAccount()
	tApp()
	tAddress()
	tQueryPreload()
	tQueryPreload2()
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
