package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


// $ export DB_USERNAME=postgres
// $ export DB_PASSWORD=postgres
// $ export DB_TABLE=postgres
// $ export DB_PORT=5432
// $ export DB_HOST=localhost

func NewDatabse() (*gorm.DB, error) {
	fmt.Println("setting up new database connection")
	dbUserName := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")

	connectString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUserName, dbTable, dbPassword)

	db, err := gorm.Open("postgres", connectString)
	if err != nil {
		return db, err
	}
	if err := db.DB().Ping(); err != nil {
		return db, err
	}
	return db, nil

}
