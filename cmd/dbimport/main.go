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
	db.Create(&bc)

	bc = m.Blockchain{Name: "Bitcoin", Description: "Bitcoin", Confirmation: 6}
	db.Create(&bc)
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

	db.Create(&asset)
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
	db.Create(&jp)
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
	db.Create(&company)
}

func tAccount() {
	var db = m.GetDB()

	// Account
	account := m.Account{
		CompanyID:       1,
		Name:            "xuyang",
		LastName:        "XU",
		FirstName:       "YANG",
		Email:           "xuyang@nbltrust.com",
		EmailEnable:     true,
		EmailVerified:   false,
		Phone:           "13912300000",
		PhoneEnabled:    true,
		PhoneVerified:   false,
		AuthKey:         "139123000001391230000013912300000",
		AuthKeyEnabled:  true,
		AuthKeyVerified: true,
		PasswordHash:    "139123000001391230000013912300000139123000001391230000013912300000139123000001391230000013912300000",
		Status:          "NORMAL",
		Type:            "ADMIN",
		Disable:         false,
	}
	db.Create(&account)
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
	db.Create(&app)
}

func tAddress() {
	var db = m.GetDB()

	address := m.Address{
		AppID:   1,
		AssetID: 1,
		Address: "0xa48d73341885e6bce0252cb05b31a8a00720cdb2",
		Status:  "NORMAL",
	}
	db.Create(&address)
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
