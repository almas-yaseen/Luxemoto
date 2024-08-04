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
		admin.GET("/PremiumCars", handlers.PremiumCars(db))
		admin.GET("/enquiry", handlers.Enquiry(db))
		admin.POST("/edit_enquiry/:id", handlers.EditEnquiry(db))
		admin.POST("/delete_enquiry/:id", handlers.DeleteEnquiry(db))
		admin.POST("/add_customers", handlers.AddCustomer(db))
		admin.GET("/MiniCars", handlers.MiniCars(db))
		admin.GET("/edit_page_premium/:id", handlers.EditPage(db))
		admin.GET("/edit_page_mini/:id", handlers.EditPageMini(db))
		admin.GET("/brand_view", handlers.BrandView(db))
		admin.POST("/brand_edit/:id", handlers.BrandEdit(db))
		admin.POST("/logout", handlers.Logout)
		admin.POST("/add_product", handlers.AddProduct(db))
		admin.POST("/edit_page_premium/:id", handlers.EditCar(db))
		admin.POST("/delete_brand/:id", handlers.BrandDelete(db))
		admin.POST("/edit_page_mini/:id", handlers.EditCarMini(db))
		admin.POST("/premium_cars/:id", handlers.PremiumCarsDelete(db))
		admin.POST("/mini_cars/:id", handlers.MiniCarsDelete(db))
		admin.POST("/add_brand", handlers.AddBrand(db))

	}

	return admin

}
