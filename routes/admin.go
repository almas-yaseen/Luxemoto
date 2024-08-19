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

		admin.POST("/logout", handlers.Logout)
		admin.GET("/get_choices", handlers.GetChoices)
		admin.GET("/dashboard", handlers.AdminDashboard(db))

		//product
		admin.GET("/product", handlers.ProductPage(db))
		admin.POST("/add_product", handlers.AddProduct(db))

		//premium cars
		admin.GET("/PremiumCars", handlers.PremiumCars(db))
		admin.POST("/premium_cars/:id", handlers.PremiumCarsDelete(db))
		admin.GET("/edit_page_premium/:id", handlers.EditPage(db))
		admin.POST("/edit_page_premium/:id", handlers.EditCar(db))
		//enquiry

		admin.GET("/enquiry", handlers.Enquiry(db))
		admin.POST("/edit_enquiry/:id", handlers.EditEnquiry(db))
		admin.POST("/delete_enquiry/:id", handlers.DeleteEnquiry(db))
		admin.POST("/add_customers", handlers.AddCustomer(db))

		//minicars

		admin.GET("/MiniCars", handlers.MiniCars(db))
		admin.GET("/edit_page_mini/:id", handlers.EditPageMini(db))
		admin.POST("/edit_page_mini/:id", handlers.EditCarMini(db))
		admin.POST("/mini_cars/:id", handlers.MiniCarsDelete(db))

		//brand

		admin.GET("/brand_view", handlers.BrandView(db))
		admin.POST("/brand_edit/:id", handlers.BrandEdit(db))
		admin.POST("/delete_brand/:id", handlers.BrandDelete(db))
		admin.POST("/add_brand", handlers.AddBrand(db))

		// Youtube

		admin.GET("/youtube_page", handlers.YoutubePage(db))
		admin.POST("/youtube_page_edit/:id", handlers.YoutubePageEdit(db))
		admin.POST("youtube_page_delete/:id", handlers.YoutubePageDelete(db))

		// Profile setup

		admin.GET("/profile", handlers.Profile(db))
		admin.POST("/change_password", handlers.ChangePassword(db))

		//gallery
		admin.GET("/gallery", handlers.Gallery(db))
		admin.POST("/add_customer_image", handlers.AddCustomerImage(db))
		admin.POST("/edit_customer_image/:id", handlers.EditCustomerImage(db))
		admin.POST("/delete_customer_images/:id", handlers.DeleteCustomerImage(db))

	}

	return admin

}
