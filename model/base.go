package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// init postgres module
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// INITFromViper ...
func INITFromViper() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
}

// InitDB ...
func InitDB(host, port, username, password, dbname string) {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, dbname, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// db.LogMode(true)
	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	// db.DropTableIfExists(&Blockchain{}, &Asset{}, &Company{}, &Account{}, &App{}, &Jadepool{}, &Order{}, &CybOrder{}, &JPOrder{}, &Event{}, &ExEvent{}, &Balance{}, &Accounting{}, &GeoAddress{}, &Address{}, &BigAsset{}, &Black{}, &Easy{})
	db.AutoMigrate(
		&JPOrder{},
		&Address{},
	)
}

//GetDB ...
func GetDB() *gorm.DB {
	return db
}
