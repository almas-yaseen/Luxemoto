package routes

import (
	"ginapp/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(myapp *gin.RouterGroup, db *gorm.DB) {
	myapp.GET("/get_premium_cars_all", handlers.GetPremiumCarsAll(db))
	myapp.GET("/get_mini_cars_all", handlers.GetMiniCarsAll(db))
	myapp.GET("/get_stock_cars_all", handlers.GetCarAll(db))
	myapp.GET("/get_specific_vehicle/:id", handlers.GetSpecificVehicle(db))
	myapp.GET("/get_choices", handlers.GetChoices)
	myapp.GET("/get_customer_images", handlers.CustomerImages(db))
	myapp.GET("/get_stock_car_all", handlers.GetStockCarAll(db))
	myapp.GET("/get_filter_types", handlers.GetFilterTypes(db))
	myapp.GET("/get_youtube_links", handlers.Get_YoutubeLink(db))
}
