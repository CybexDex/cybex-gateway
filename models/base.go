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

	// username := "postgres"      // os.Getenv("db_user")
	// password := "your_password" // os.Getenv("db_pass")
	// dbName := "demo"            // os.Getenv("db_name")
	// dbHost := "47.100.239.204"  // os.Getenv("db_host")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn

	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	// FOR TEST USE ONLY !!!
	// db.DropTableIfExists(&Blockchain{}, &Asset{}, &Company{}, &Account{}, &App{}, &Jadepool{}, &Order{}, &ExOrder{}, &Event{}, &ExEvent{}, &Balance{}, &Accounting{}, &GeoAddress{}, &Address{})

	db.AutoMigrate(&Blockchain{}, &Asset{}, &Company{}, &Account{}, &App{}, &Jadepool{}, &Order{}, &ExOrder{}, &Event{}, &ExEvent{}, &Balance{}, &Accounting{}, &GeoAddress{}, &Address{})
}

//GetDB ...
func GetDB() *gorm.DB {
	return db
}
