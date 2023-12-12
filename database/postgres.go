package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pc01pc013/task-management/database/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

func init() {
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	mode := os.Getenv("MODE")
	gormConfig := &gorm.Config{}
	if mode == "Test" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), gormConfig)

	if err != nil {
		log.Fatalf("Error opening database (gorm): %q", err)
	}

	gormDB.AutoMigrate(&entities.User{}, &entities.Label{}, &entities.Task{}, &entities.Userinfo_Google{})
}

func GetDBInstance() *gorm.DB {
	if gormDB == nil {
		log.Fatalf("Error getDBInstance")
	} else {
		return gormDB
	}
	return nil
}

func Close() error {
	sqlDB, _ := gormDB.DB()
	return sqlDB.Close()
}
