package database

import (
	"auth_api_with_Go/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

var Instance *gorm.DB
var dbError error

func NewConnection(config Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.User, config.Password, config.Host, config.Port, config.DBName,
    )

    Instance, dbError = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if dbError != nil {
        log.Fatal(dbError)
        return nil, dbError
    }
	log.Println("Connected to Database!")
    return Instance, nil
}

func Migrate() {
	Instance.AutoMigrate(&model.User{})
	log.Println("Database Migration Completed!")
}