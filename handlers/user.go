package handlers

import (
	"ginapp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPremiumCarsAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cars []domain.Vehicle

		if err := db.Order("created_at desc").Limit(6).Preload("Brand").Preload("Images").Where("vehicle_type = ? ", "Premium").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the cars"})
		}

		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Year         int    `json:"int"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			RegNo        string `json:"reg_no"`
			Status       string `json:"status"`
			Price        int    `json:"int"`
			Image        string `json:"image"`
		}

		var result []CarWithImage

		for _, car := range cars {
			var image string
			if len(car.Images) > 0 {

				image = car.Images[0].Path
			}

			CarWithImage := CarWithImage{

				ID:           car.ID,
				Brand:        car.Brand.Name,
				Model:        car.Model,
				Year:         car.Year,
				Color:        car.Color,
				Variant:      car.Variant,
				Kms:          car.Kms,
				Ownership:    car.Ownership,
				FuelType:     car.FuelType,
				CarType:      car.CarType,
				Transmission: car.Transmission,
				RegNo:        car.RegNo,
				Status:       car.Status,
				Price:        car.Price,
				Image:        image,
			}
			result = append(result, CarWithImage)
		}

		c.JSON(http.StatusOK, gin.H{"premium_cars": result})

	}
}

func GetMiniCarsAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
