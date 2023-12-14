package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pc01pc013/task-management/database/entities"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

func init() {
	mode := os.Getenv("MODE")
	gormConfig := &gorm.Config{}
	if mode == "Test" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var sqlDB *sql.DB
	var err error

	switch os.Getenv("DATABASE_DRIVER") {
	case "postgres":
		sqlDB, err = sql.Open(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}

		gormDB, err = gorm.Open(postgres.New(postgres.Config{
			Conn:                 sqlDB,
			PreferSimpleProtocol: true,
		}), gormConfig)

	case "sqlserver":
		gormDB, err = gorm.Open(sqlserver.Open(os.Getenv("DATABASE_URL")), gormConfig)
	default:
		log.Fatalf("Error DATABASE_DRIVER")
	}

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
