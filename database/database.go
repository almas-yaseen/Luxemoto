package database

import (
	"fmt"
	"ginapp/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	fmt.Println("here is the sqlinfo", psqlInfo)
	db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if dberr != nil {
		return nil, fmt.Errorf("Failed to connect the database")

	}
	DB = db

	return DB, nil

}
