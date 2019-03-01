package main

import (
	"fmt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// init config
	utils.InitConfig()

	// init db
	dbHost := viper.GetString("database.host")
	username := viper.GetString("database.user")
	password := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

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
