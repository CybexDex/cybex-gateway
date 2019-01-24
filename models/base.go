package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"

	"github.com/joho/godotenv"
	// init postgres module
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	// username := "alexxu"
	// password := "postgres"
	// dbName := "xuyang"
	// dbHost := "127.0.0.1"

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	// db.LogMode(true)
	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	db.DropTableIfExists(&Blockchain{}, &Asset{}, &Company{}, &Account{}, &App{}, &Jadepool{}, &Order{}, &ExOrder{}, &Event{}, &ExEvent{}, &Balance{}, &Accounting{}, &GeoAddress{}, &Address{}, &CybToken{})

	db.AutoMigrate(&Blockchain{}, &Asset{}, &Company{}, &Account{}, &App{}, &Jadepool{}, &Order{}, &ExOrder{}, &Event{}, &ExEvent{}, &Balance{}, &Accounting{}, &GeoAddress{}, &Address{}, &CybToken{})
}

//GetDB ...
func GetDB() *gorm.DB {
	return db
}
