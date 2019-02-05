package main

import (
	"fmt"
	"os"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	// username := "alexxu"
	// password := "postgres"
	// dbName := "xuyang"
	// dbHost := "localhost"

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	conn.LogMode(true)

	conn.DropTableIfExists(&m.Blockchain{}, &m.Asset{}, &m.Company{}, &m.Account{}, &m.App{},
		&m.Jadepool{}, &m.Order{}, &m.CybOrder{}, &m.JPOrder{}, &m.Event{},
		&m.ExEvent{}, &m.Balance{}, &m.Accounting{}, &m.GeoAddress{}, &m.Address{}, &m.CybToken{})

	conn.AutoMigrate(&m.Blockchain{}, &m.Asset{}, &m.Company{}, &m.Account{}, &m.App{},
		&m.Jadepool{}, &m.Order{}, &m.CybOrder{}, &m.JPOrder{}, &m.Event{},
		&m.ExEvent{}, &m.Balance{}, &m.Accounting{}, &m.GeoAddress{}, &m.Address{}, &m.CybToken{})
}
