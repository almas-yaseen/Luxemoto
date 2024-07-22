package routes

import (
	"ginapp/handlers"
	"ginapp/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(admin *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	admin.GET("/login", handlers.Adminlogin)
	admin.POST("/login", handlers.AdminLogin(db))
	admin.Use(middleware.AdminAuthMiddleware(db))
	{
		admin.GET("/get_choices", handlers.GetChoices)
		admin.GET("/dashboard", handlers.AdminDashboard)
		admin.GET("/product", handlers.ProductPage(db))
		admin.GET("/PremiumCars", handlers.PremiumCars)
		admin.POST("/logout", handlers.Logout)
		admin.POST("/add_product", handlers.AddProduct(db))

	}

	return admin

}
