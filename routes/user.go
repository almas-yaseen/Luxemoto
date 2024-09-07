package routes

import (
	"ginapp/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(myapp *gin.RouterGroup, db *gorm.DB) {
	myapp.GET("/get_premium_cars_all", handlers.GetPremiumCarsAll(db))
	myapp.GET("/get_mini_cars_all", handlers.GetMiniCarsAll(db))
	myapp.GET("/get_stock_cars_all/*type", handlers.GetCarAll(db))
	myapp.GET("/get_specific_vehicle/:id", handlers.GetSpecificVehicle(db))
	myapp.GET("/get_choices", handlers.GetChoices)
	myapp.GET("/get_latest_delivery", handlers.CustomerImages(db))
	myapp.GET("/get_filter_types", handlers.GetFilterTypes(db))
	myapp.GET("/get_youtube_links", handlers.Get_YoutubeLink(db))
	myapp.GET("/get_all_delivery", handlers.GetAllDelivery(db))
	myapp.GET("/get_banner_details", handlers.GetBannerDetails(db))
}
