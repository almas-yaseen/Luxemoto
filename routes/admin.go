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
		admin.GET("/dashboard", handlers.AdminDashboard)
		admin.GET("/product", handlers.ProductPage)
		admin.POST("/logout", handlers.Logout)

	}

	return admin

}
