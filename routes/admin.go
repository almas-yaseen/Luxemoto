package routes

import (
	"ginapp/handlers"
	"ginapp/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(admin *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	admin.GET("/login", handlers.Adminlogin)
	admin.Use(middleware.AdminAuthMiddleware(db))
	admin.GET("/dashboard", handlers.AdminDashboard)
	admin.POST("/login", handlers.AdminLogin(db))

	return admin

}
