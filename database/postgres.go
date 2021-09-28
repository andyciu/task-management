package database

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func init() {
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error opening database (gorm): %q", err)
	}
}

func GetDBInstance() *gorm.DB {
	if gormDB == nil {
		log.Fatalf("Error getDBInstance")
	} else {
		return gormDB
	}
	return nil
}
