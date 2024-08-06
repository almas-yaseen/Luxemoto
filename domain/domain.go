package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	CarStatusSold      = "sold"
	CarStatusAvailable = "available"
)

const (
	CarTypeSedan       = "sedan"
	CarTypeHatchback   = "hatchback"
	CarTypeSuv         = "suv"
	CarTypeMpv         = "mpv"
	CarTypeCompactSuv  = "compact suv"
	CarTypeConvertible = "convertible"
	CarTypeSport       = "sport"

	// fuel type

	FuelTypePetrol   = "petrol"
	FuelTypeDiesel   = "diesel"
	FuelTypeHybrid   = "hybrid"
	FuelTypeCng      = "cng"
	FuelTypeElectric = "electric"

	// category
	CarCategoryPremium = "Premium"
	CarCategoryMini    = "Mini"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"unique;not null"`
	IsAdmin   bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Vehicle struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"json:"id"`
	Model          string    `json:"model"`
	BrandID        uint      `json:"brand_id"`
	Brand          Brand     `gorm:"foreignKey:BrandID" json:"brand"`
	Year           int       `json:"year"`
	Category       string    `json:"category"`
	Color          string    `json:"color"`
	CarType        string    `json:"car_type"`
	Images         []Image   `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE"json:"images"`
	FuelType       string    `json:"fuel_type"`
	Variant        string    `json:"variant"`
	Kms            int       `json:"kms"`
	Ownership      int       `json:"ownership"`
	Vehicle_type   string    `json:"vehicle_type"`
	Transmission   string    `json:"transmission"`
	RegNo          string    `json:"regno"`
	Status         string    `json:"status"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Engine_size    string    `json:"engine_size"`
	Insurance_date string    `json:"insurance_dating"`
	Location       string    `json:"location"`
}

type Image struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CarID uint   `json:"car_id"`
	Path  string `json:"path"`
}

type Brand struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"json:"id"`
	Name      string    `gorm:"unique;" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Enquiry struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"json:"id"`
	CustomerName string `json:"customer_name"`
	Phone        int    `gorm:"not null"`
	DesiredPrice int    `gorm:"not null"`
	DesiredCars  string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CustomerImage struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uint   `json:"id"`
	Path      string `json:"path"`
}

type YoutubeLink struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uint      `json:"id"`
	VideoLink string    `json:"video_links"`
}
