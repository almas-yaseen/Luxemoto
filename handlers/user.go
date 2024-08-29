package handlers

import (
	"fmt"
	"ginapp/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBannerDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cars []domain.Vehicle
		if err := db.Order("created_at desc").Limit(5).Preload("Brand").Preload("Images").Where("vehicle_type= ?", "Premium").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tha database"})
			return

		}

		// Create a structure to hold the response data
		type CarDetail struct {
			BannerImage string `json:"bannerImage"`
			Brand       string `json:"brand"`
			Id          int    `json:"id"`
			Cartype     string `json:"car_type"`
			VehicleType string `json:"vehicle_type"`
			Year        int    `json:"year"`
			Model       string `json:"model"`
			Variant     string `json:"variant"`
			Price       int    `json:"price"`
			Color       string `json:"color"`
		}

		var carDetails []CarDetail

		for _, car := range cars {

			carDetail := CarDetail{
				BannerImage: car.BannerImage,
				Model:       car.Model,
				Variant:     car.Variant,
				Price:       car.Price,
				Color:       car.Color,
				VehicleType: car.Vehicle_type,
				Cartype:     car.CarType,
				Brand:       car.Brand.Name,
				Year:        int(car.Year),
				Id:          int(car.ID),
			}
			carDetails = append(carDetails, carDetail)
			fmt.Println("Car details:", carDetail)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "carDetails": carDetails})
	}
}
func Get_YoutubeLink(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			youtubelinks []domain.YoutubeLink
			page         int
			limit        int
			offset       int
			totalCount   int64
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset = (page - 1) * limit

		if err := db.Model(&domain.YoutubeLink{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the count"})
			return
		}
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&youtubelinks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"youtubelinks": youtubelinks})
	}
}

func GetFilterTypes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filtertypes struct {
			CarTypes  []string       `json:"car_types"`
			FuelTypes []string       `json:"fuel_types"`
			Brands    []domain.Brand `json:"brands"`
		}

		var brands []domain.Brand // getting the values

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brands"}) // storing the val
			return
		}

		filtertypes.Brands = brands

		filtertypes.CarTypes = []string{
			domain.CarCategoryMini,
			domain.CarCategoryPremium,
			domain.CarTypeCompactSuv,
			domain.CarTypeConvertible,
			domain.CarTypeHatchback,
			domain.CarTypeMpv,
		}

		filtertypes.FuelTypes = []string{
			domain.FuelTypeCng,
			domain.FuelTypeDiesel,
			domain.FuelTypeElectric,
			domain.FuelTypeHybrid,
		}

		c.JSON(http.StatusOK, filtertypes)
	}
}

func GetStockCarAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			cars   []domain.Vehicle
			limit  int
			offset int
			count  int64
			page   int
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset = (page - 1) * limit

		BrandIDStr := c.Query("brand_id")
		carType := c.Query("car_type")
		fuelType := c.Query("fuel_type")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")
		query := db.Model(&domain.Vehicle{})

		if BrandIDStr != "" {
			brandID, err := strconv.ParseUint(BrandIDStr, 10, 64)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brandID format"})
				return
			}
			query = query.Where("brand_id=?", brandID)

		}

		if carType != "" {
			query = query.Where("car_type = ?", carType)
		}

		if fuelType != "" {
			query = query.Where("fuel_type= ?", fuelType)
		}
		if minPrice != "" {
			minPriceFloat, err := strconv.ParseFloat(minPrice, 64)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price format"})
				return
			}

			query = query.Where("min_price >= ?", minPriceFloat)
		}

		if maxPrice != "" {
			maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
				return
			}
			query = query.Where("price <= ?", maxPriceFloat)
		}

		if err := query.Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count the cars"})
			return
		}

		if err := query.Order("created_at desc").Preload("Brand").Preload("Images").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to  find the cars"})
			return
		}

		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Status       string `json:"status"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
			Variant      string `json:"variant"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			Price        int    `json:"price"`
			Image        string `json:"image"`
		}
		var result []CarWithImage

		for _, cars := range cars {

			var image string
			if len(cars.Images) > 0 {
				image = cars.Images[0].Path

			}

			CarWithImage := CarWithImage{
				ID:           cars.ID,
				Brand:        cars.Brand.Name,
				Model:        cars.Model,
				Year:         cars.Year,
				Color:        cars.Color,
				CarType:      cars.CarType,
				FuelType:     cars.FuelType,
				Variant:      cars.Variant,
				Kms:          cars.Kms,
				Ownership:    cars.Ownership,
				Transmission: cars.Transmission,
				Status:       cars.Status,
				Price:        cars.Price,
				Image:        image,
			}

			result = append(result, CarWithImage)

		}
		c.JSON(http.StatusOK, gin.H{"status": "200", "all_stocks": result})

	}
}

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
			Year         int    `json:"year"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			RegNo        string `json:"reg_no"`
			Status       string `json:"status"`
			Price        int    `json:"price"`
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

		var cars []domain.Vehicle

		if err := db.Order("created_at desc").Limit(6).Preload("Brand").Preload("Images").Where("vehicle_type = ? ", "Mini").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brand"})
			return
		}
		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			RegNo        string `json:"reg_no"`
			Status       string `json:"status"`
			Price        int    `json:"price"`
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
		c.JSON(http.StatusOK, gin.H{"mini_cars": result})

	}
}

func GetCarAll(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cars []domain.Vehicle

		if err := db.Order("created_at desc").Preload("Brand").Preload("Images").Find(&cars).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brand"})
			return
		}

		type CarWithImage struct {
			ID           uint   `json:"id"`
			Brand        string `json:"brand"`
			Model        string `json:"model"`
			Year         int    `json:"year"`
			Color        string `json:"color"`
			Variant      string `json:"variant"`
			CarType      string `json:"car_type"`
			FuelType     string `json:"fuel_type"`
			Kms          int    `json:"kms"`
			Ownership    int    `json:"ownership"`
			Transmission string `json:"transmission"`
			RegNo        string `json:"reg_no"`
			Status       string `json:"status"`
			Price        int    `json:"price"`
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

		ctx.JSON(http.StatusOK, gin.H{"all_cars": result})

	}
}

func GetSpecificVehicle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		fmt.Println("here is the id")

		var cars domain.Vehicle

		if err := db.Preload("Brand").Preload("Images").First(&cars, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to  find the id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"vehicle": cars})

	}
}

func CustomerImages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			CustomerImages []domain.CustomerImage
			page           int
			limit          int
			totalCount     int64
			offset         int
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit

		if err := db.Model(&domain.CustomerImage{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to  find the image"})
			return
		}

		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&CustomerImages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the customer images"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"images": CustomerImages, "TotalCount": totalCount})

	}
}

func GetAllDelivery(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			CustomerImages []domain.CustomerImage
			page           int
			limit          int
			totalCount     int64
			offset         int
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit

		if err := db.Model(&domain.CustomerImage{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to  find the image"})
			return
		}

		if err := db.Limit(limit).Offset(offset).Find(&CustomerImages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the customer images"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"CustomerImages": CustomerImages, "TotalCount": totalCount})

	}
}
