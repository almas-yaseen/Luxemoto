package routes

import (
	"ginapp/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(myapp *gin.RouterGroup, db *gorm.DB) {
	myapp.GET("/get_premium_cars_all", handlers.GetPremiumCarsAll(db))
	myapp.GET("/get_mini_cars_all", handlers.GetMiniCarsAll(db))
}
