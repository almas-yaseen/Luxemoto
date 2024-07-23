package database

import (
	"fmt"
	"ginapp/config"
	"ginapp/domain"

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
	fmt.Println("here is DB", DB)

	DB.AutoMigrate(&domain.User{})
	DB.AutoMigrate(&domain.Brand{})
	DB.AutoMigrate(&domain.Vehicle{})
	DB.AutoMigrate(&domain.Image{})

	return DB, nil

}
