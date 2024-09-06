package handlers

import (
	"encoding/json"
	"fmt"
	"ginapp/domain"
	"ginapp/services"
	"ginapp/utils"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddYoutube(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var youtube domain.YoutubeLink
		youtubelink := c.PostForm("youtube_link") // Corrected key here
		fmt.Println("here is the youtubelink", youtubelink)

		youtube.VideoLink = youtubelink
		if err := db.Create(&youtube).Error; err != nil {
			// Corrected error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add the YouTube video"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/youtube_page")
	}
}

func AdminDashboard(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var (
			AboveFiftyLakh    []domain.Vehicle
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
			YoutubeVideoCount int64
			brands            []domain.Brand
			brandCounts       = make(map[string]int64)
		)
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		if err := db.Preload("Brand").Model(&domain.Vehicle{}).Where("price > ?", 5000000).Find(&AboveFiftyLakh).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the above fifty lakh customer"})
			return
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the brands"})
			return
		}

		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}
		if err := db.Model(&domain.YoutubeLink{}).Count(&YoutubeVideoCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count the youtube videos"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		for _, brand := range brands {
			var count int64
			if err := db.Model(&domain.Vehicle{}).Where("brand_id = ?", brand.ID).Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count cars per brand"})
				return
			}
			brandCounts[brand.Name] = count
		}

		// Convert brandCounts to JSON
		brandCountsJSON, err := json.Marshal(brandCounts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal brand counts"})
			return
		}

		// Render the dashboard with the retrieved data
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"abovefifty":        AboveFiftyLakh,
			"brands":            brands,
			"brandCounts":       brandCounts,
			"brandCountsJ":      template.JS(brandCountsJSON), // Pass brandCounts as JSON
			"CurrentPath":       c.Request.URL.Path,
			"youtubecount":      YoutubeVideoCount,
		})
	}
}

func DeleteCustomerImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var customerimage domain.CustomerImage

		id := c.Param("id")

		if err := db.First(&customerimage, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete the customer image"})
			return
		}

		if err := os.Remove(customerimage.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed remove the customer image"})
			return
		}

		if err := db.Delete(&customerimage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete the image"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/gallery")

	}
}

func EditCustomerImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customerimage domain.CustomerImage
		id := c.Param("id")
		if err := db.First(&customerimage, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to  find the id"})
			return
		}
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to  find the image"})
			return
		}

		newimagepath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
		if err := c.SaveUploadedFile(file, newimagepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
			return
		}

		customerimage.Path = newimagepath
		if err := db.Save(&customerimage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the image"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/gallery")

	}
}

func AddCustomerImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var customerimages domain.CustomerImage

		CustomerImage, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the image"})
			return
		}
		fmt.Println("here is teh customer image", CustomerImage)
		customerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), CustomerImage.Filename)) //creating  the file and the details
		if err := c.SaveUploadedFile(CustomerImage, customerImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
			return
		}

		customerimages.Path = customerImagePath

		if err := db.Create(&customerimages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/gallery")

	}
}

func Gallery(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var (
			gallery           []domain.CustomerImage
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
			page              int
			limit             int
			offset            int
			totalCount        int64
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "1"))

		offset = (page - 1) * limit

		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&gallery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer image failed to get"})
			return
		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.CustomerImage{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)

		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(http.StatusOK, "gallery.html", gin.H{"enquries": gallery, "CurrentPath": c.Request.URL.Path, "enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal": PremiumCarsTotal,
			"minicarstotal":    MiniCarsTotal,
			"Page":             page,
			"Limit":            limit,
			"totalPages":       totalPages,
			"totalCount":       totalCount,
			"totalcars":        TotalCars})

	}
}

func ChangePassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		email := c.PostForm("email")
		fmt.Println("here is the email", email)
		current_password := c.PostForm("current_password")
		fmt.Println("here is the current password", current_password)
		new_password := c.PostForm("password1")
		confirm_password := c.PostForm("password2")

		var user domain.User

		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			fmt.Println("Error finding user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the user"})
			return
		}
		fmt.Println("Stored Hashed Password:", user.Password)

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current_password)); err != nil {
			fmt.Println("Input Current Password:", current_password)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Incorrect current password",
				"details": err.Error(), // Include the actual error message for debugging
			})
			return
		}

		if new_password != confirm_password {
			c.JSON(http.StatusBadRequest, gin.H{"error": "new password or confirm password not correct "})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash the password"})
			return
		}

		user.Password = string(hashedPassword)
		fmt.Println("here is the user.password", user.Password)

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the user"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/profile")

	}
}

func Profile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var user []domain.User

		if err := db.Find(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the user"})
			return
		}
		c.HTML(http.StatusOK, "profile.html", gin.H{"user": user, "CurrentPath": c.Request.URL.Path})

	}
}

func YoutubePageDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id")
		var youtube domain.YoutubeLink

		if err := db.First(&youtube, id).Find(&youtube).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the youtube link"})
			return
		}

		if err := db.Delete(&youtube).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"deleted": "successfully"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/youtube_page")
	}
}

func YoutubePageEdit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id")
		var youtubelink domain.YoutubeLink

		if err := db.First(&youtubelink, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the youtube id"})
			return
		}

		youtubelink.VideoLink = c.PostForm("link")

		if err := db.Save(&youtubelink).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the car"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/youtube_page")

	}
}

func YoutubePage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var (
			youtubelink       []domain.YoutubeLink
			page              int
			limit             int
			offset            int
			totalCount        int64
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)
		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))

		if page < 1 {
			page = 1
		}

		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit

		if err := db.Model(&domain.YoutubeLink{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInsufficientStorage, gin.H{"error": "failed to the count"})
		}

		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&youtubelink).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed find the youtube link"})
		}
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)
		fmt.Println("here is the youtubelink", youtubelink)

		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(http.StatusOK, "youtube.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"Pages":             pages,
			"Limit":             limit,
			"totalCount":        totalCount,
			"totalPages":        totalPages,
			"Page":              page,
			"CurrentPath":       c.Request.URL.Path,
			"enquries":          youtubelink,
		})

	}
}

func DeleteEnquiry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		id := c.Param("id")
		var enquiry domain.Enquiry
		fmt.Println("here is the address value", &enquiry)
		if err := db.First(&enquiry, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the id"})
			return
		}

		if err := db.Delete(&enquiry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the database"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/enquiry")

	}
}

func EditEnquiry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id")

		var enquiry domain.Enquiry

		if err := db.First(&enquiry, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the id"})
			return
		}

		enquiry.CustomerName = c.PostForm("customer_name")
		enquiry.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		enquiry.DesiredPrice, _ = strconv.Atoi(c.PostForm("price"))
		enquiry.DesiredCars = c.PostForm("cars")

		if err := db.Order("created_at desc").Save(&enquiry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the database"})
			return
		}
		fmt.Println("here is the saved enquiry", enquiry)

		c.Redirect(http.StatusSeeOther, "/admin/enquiry")

	}
}

func AddCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var Customer domain.Enquiry

		Customer.CustomerName = c.PostForm("customer_name")
		Customer.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		Customer.DesiredPrice, _ = strconv.Atoi(c.PostForm("price"))
		Customer.DesiredCars = c.PostForm("cars")
		if err := db.Create(&Customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create the car"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/enquiry")

	}
}

func Enquiry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var (
			enquires          []domain.Enquiry
			page              int
			limit             int
			offset            int
			totalCount        int64
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "1"))
		offset = (page - 1) * limit

		if err := db.Model(&domain.Enquiry{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the count"})
		}
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&enquires).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the data"})
			return
		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}
		fmt.Println("here is the totalcount", totalCount)
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)

		for i := range pages {
			pages[i] = i + 1

		}

		fmt.Println("HERE IS THE QNUER", enquires)
		c.HTML(http.StatusOK, "enquiry.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"enquries":          enquires,
			"Page":              page,
			"CurrentPath":       c.Request.URL.Path,
			"Pages":             pages,
			"Limit":             limit,
			"totalCount":        totalCount,
			"totalPages":        totalPages,
		})

	}
}
func AddBrand(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var brands domain.Brand
		brands.Name = c.PostForm("brand_name")

		if err := db.Create(&brands).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/brand_view")

		}
		c.Redirect(http.StatusSeeOther, "/admin/brand_view")

	}
}

func MiniCarsDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var cars []domain.Vehicle

		id := c.Param("id")

		if err := db.First(&cars, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the cars"})
			return
		}

		if err := db.Delete(&cars).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/MiniCars")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/MiniCars")
	}
}

func PremiumCarsDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var cars domain.Vehicle
		id := c.Param("id")
		if err := db.First(&cars, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the id"})
			return
		}

		if err := db.Delete(&cars).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/PremiumCars") // Redirect to brands list page
			return

		}
		c.Redirect(http.StatusSeeOther, "/admin/PremiumCars")

	}
}

func BrandDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var brand domain.Brand
		id := c.Param("id")

		// Check if the brand exists
		if err := db.First(&brand, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}

		// here am
		if err := db.Model(&domain.Vehicle{}).Where("brand_id = ?", brand.ID).Update("brand_id", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update associated cars"})
			return
		}

		// Now delete the brand itself
		if err := db.Delete(&brand).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/brand_view") // Redirect to brands list page
		}

		// Success message or redirect
		c.Redirect(http.StatusSeeOther, "/admin/brand_view") // Redirect to brands list page
	}
}

func BrandEdit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var brand domain.Brand
		id := c.Param("id")
		fmt.Println("here is the id", id)

		// Validate if brand with given ID exists
		if err := db.First(&brand, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
			return
		}

		// Bind form data to get the new brand name
		if err := c.ShouldBind(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		newBrandName := c.PostForm("brand_name")

		// Update brand name
		brand.Name = newBrandName

		// Save updated brand to the database
		if err := db.Save(&brand).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/brand_view") // Redirect to brands list page
		}

		// Redirect or respond with success message
		c.Redirect(http.StatusSeeOther, "/admin/brand_view") // Redirect to brands list page
	}
}

func BrandView(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var (
			brands            []domain.Brand
			totalCount        int64
			page              int
			limit             int
			offset            int
			PremiumCarsTotal  int64
			MiniCarsTotal     int64 //newone
			TotalCars         int64 //newone
			YoutubeVideoCount int64
			EnquiredCustomers []domain.Enquiry
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))

		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit

		fmt.Println("here is the brands", brands)

		if err := db.Model(&domain.Brand{}).Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the count"})
			return
		}
		if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the brands"})

			return
		}

		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the premium cars"})
			return
		}

		if err := db.Model(&domain.YoutubeLink{}).Count(&YoutubeVideoCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the totalcount"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the mini cars"})
			return
		}
		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the total count of the cars"})
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)

		for i := range pages {

			pages[i] = i + 1
		}

		noBrands := len(brands) == 0

		c.HTML(http.StatusOK, "brand.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"brands":            brands,
			"totalCount":        totalCount,
			"Page":              page,
			"noBrands":          noBrands,
			"Limit":             limit,
			"totalPages":        totalPages,
			"youtubecount":      YoutubeVideoCount,
			"Pages":             pages,
			"CurrentPath":       c.Request.URL.Path, // Pass the current path to the template
		})

	}
}

func EditCarMini(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id")
		var car domain.Vehicle
		var brands []domain.Brand

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the id"})
			return
		}
		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brands"})
			return
		}

		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the year"})
			return
		}

		car.Year = year
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Vehicle_type = c.PostForm("vehicle_type")
		car.Status = c.PostForm("status")
		car.Price, _ = strconv.Atoi(c.PostForm("price"))
		car.CarType = c.PostForm("car_type")
		car.FuelType = c.PostForm("fuel_type")
		car.Engine_size = c.PostForm("engine_size")
		car.Insurance_date = c.PostForm("insurance_date")
		car.Location = c.PostForm("location")

		form, err := c.MultipartForm()

		if err == nil {

			files := form.File["images[]"]
			fmt.Println("here is the mini images", files)
			var newImages []domain.Image

			for _, file := range files {
				filename := filepath.Base(fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename)) // creating the base filename  with unique values
				uploadPath := filepath.Join("uploads", filename)                                                 // this tell where you want to exactly store the data
				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
					return
				}

				imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				newImages = append(newImages, domain.Image{Path: imagePath})
			}
			if len(newImages) > 0 {
				car.Images = append(car.Images, newImages...)
			}

		}

		deleteImageIDs := c.PostFormArray("delete_images")

		if len(deleteImageIDs) > 0 {
			var remainingImages []domain.Image
			for _, img := range car.Images {
				shouldDelete := false
				for _, id := range deleteImageIDs {
					if strconv.Itoa(int(img.ID)) == id {
						shouldDelete = true
						break
					}
				}
				if !shouldDelete {
					remainingImages = append(remainingImages, img)
				} else {
					if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
						fmt.Print("failed to delete the image", err)

					}
					db.Delete(&img)
				}
			}
			car.Images = remainingImages
		}

		for _, img := range car.Images {
			file, err := c.FormFile(fmt.Sprintf("replace_image_%d", img.ID))

			if err == nil {
				if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
					fmt.Println("failed to delete the image in replace", err)
				}

				uploadPath := filepath.Join("uploads", fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename))

				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
					return
				}

				img.Path = "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				db.Save(&img)
			}
		}
		brandID, err := strconv.ParseUint(c.PostForm("brand"), 10, 64)
		fmt.Println("here is the brandID", brandID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand id"})
			return
		}
		car.Brand.ID = uint(brandID)
		fmt.Println("here is the car.Brandid", car.BrandID)
		if err := db.Save(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update the car"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/MiniCars")

	}

}

func MiniCars(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		fmt.Println("here is the mini cars come on")

		var (
			cars              []domain.Vehicle
			brands            []domain.Brand
			totalCount        int64
			limit             int
			offset            int
			page              int
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit
		if err := db.Preload("Brand").Order("created_at desc").Limit(limit).Offset(offset).Where("vehicle_type=?", "Mini").Find(&cars).Error; err != nil {
			fmt.Println("jerandlnaskljd")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the brands"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find  the total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the vehicle"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type=?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the premium cars"})
			return
		}

		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the enquired customers"})
			return
		}

		for i := range cars {

			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the images"})
				return

			}
		}

		noCars := len(cars) == 0

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)

		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(200, "MiniCars.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"Cars":              cars,
			"Brands":            brands,
			"NoCars":            noCars,
			"Page":              page,
			"Pages":             pages,
			"totalPages":        totalPages,
			"CurrentPath":       c.Request.URL.Path,
			"Limit":             limit,
			"totalCount":        totalCount,
			"Nocars":            noCars,
		})

	}

}

func EditCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id")
		var car domain.Vehicle
		var brands []domain.Brand

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the id"})
			return
		}
		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brands"})
			return
		}

		car.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the year"})
			return
		}

		car.Year = year
		car.Color = c.PostForm("color")
		car.Variant = c.PostForm("variant")
		car.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		car.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		car.Transmission = c.PostForm("transmission")
		car.RegNo = c.PostForm("regno")
		car.Vehicle_type = c.PostForm("vehicle_type")
		car.Status = c.PostForm("status")
		car.Price, _ = strconv.Atoi(c.PostForm("price"))
		car.CarType = c.PostForm("car_type")
		car.FuelType = c.PostForm("fuel_type")
		car.Engine_size = c.PostForm("engine_size")
		car.Insurance_date = c.PostForm("insurance_date")
		car.Location = c.PostForm("location")

		bannerImage, err := c.FormFile("bannerimage")
		if err == nil {
			// Delete the old banner image if it exists
			if car.BannerImage != "" {
				if err := deleteFile(strings.TrimPrefix(car.BannerImage, "/")); err != nil {
					fmt.Println("Failed to delete the old banner image:", err)
				}
			}

			// Upload new banner image
			bannerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", car.ID, bannerImage.Filename))
			if err := c.SaveUploadedFile(bannerImage, bannerImagePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the banner image"})
				return
			}
			bannerImagePath = "/" + strings.ReplaceAll(bannerImagePath, "\\", "/")
			car.BannerImage = bannerImagePath
		}

		form, err := c.MultipartForm()

		if err == nil {

			files := form.File["images[]"]
			fmt.Println("here is the imagesxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", files)

			var newImages []domain.Image

			for _, file := range files {
				filename := filepath.Base(fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename)) // creating the base filename  with unique values
				uploadPath := filepath.Join("uploads", filename)                                                 // this tell where you want to exactly store the data
				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
					return
				}

				imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				newImages = append(newImages, domain.Image{Path: imagePath})
			}
			if len(newImages) > 0 {
				car.Images = append(car.Images, newImages...)
			}

		}

		deleteImageIDs := c.PostFormArray("delete_images")

		if len(deleteImageIDs) > 0 {
			var remainingImages []domain.Image
			for _, img := range car.Images {
				shouldDelete := false
				for _, id := range deleteImageIDs {
					if strconv.Itoa(int(img.ID)) == id {
						shouldDelete = true
						break
					}
				}
				if !shouldDelete {
					remainingImages = append(remainingImages, img)
				} else {
					if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
						fmt.Print("failed to delete the image", err)

					}
					db.Delete(&img)
				}
			}
			car.Images = remainingImages
		}

		for _, img := range car.Images {
			file, err := c.FormFile(fmt.Sprintf("replace_image_%d", img.ID))

			if err == nil {
				if err := deleteFile(strings.TrimPrefix(img.Path, "/")); err != nil {
					fmt.Println("failed to delete the image in replace", err)
				}

				uploadPath := filepath.Join("uploads", fmt.Sprintf("%d_%d_%s", car.ID, time.Now().UnixNano(), file.Filename))

				if err := c.SaveUploadedFile(file, uploadPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the image"})
					return
				}

				img.Path = "/" + strings.ReplaceAll(uploadPath, "\\", "/")
				db.Save(&img)
			}
		}
		brandID, err := strconv.ParseUint(c.PostForm("brand"), 10, 64)
		fmt.Println("here is the brandID")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand id"})
			return
		}
		car.Brand.ID = uint(brandID)
		if err := db.Save(&car).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update the car"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/PremiumCars")

	}

}

func deleteFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("filepath is empty")
	}
	err := os.Remove(filepath)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Files does not exist %s", filepath)
		}

		return err

	}
	return nil

}

func EditPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		id := c.Param("id") // fetching the id from the url
		var (
			car               domain.Vehicle
			brands            []domain.Brand
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the id"})
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the car"})
			return
		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		CarTypes := []string{
			domain.CarTypeCompactSuv,
			domain.CarTypeConvertible,
			domain.CarTypeMpv,
			domain.CarTypeHatchback,
			domain.CarTypeSedan,
			domain.CarTypeSport,
		}
		FuelTypes := []string{
			domain.FuelTypeCng,
			domain.FuelTypeDiesel,
			domain.FuelTypeElectric,
			domain.FuelTypePetrol,
			domain.FuelTypeHybrid,
		}
		VehicleTypes := []string{
			domain.CarCategoryMini,
			domain.CarCategoryPremium,
		}

		c.HTML(http.StatusOK, "edit_premium.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"fuel_types":        FuelTypes,
			"car_types":         CarTypes,
			"Car":               car,
			"Brands":            brands,
			"vehicle_types":     VehicleTypes,
		})

	}
}
func EditPageMini(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") // fetching the id from the url
		var (
			car               domain.Vehicle
			brands            []domain.Brand
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the id"})
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the car"})
			return
		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		CarTypes := []string{
			domain.CarTypeCompactSuv,
			domain.CarTypeConvertible,
			domain.CarTypeMpv,
			domain.CarTypeHatchback,
			domain.CarTypeSedan,
			domain.CarTypeSport,
		}
		FuelTypes := []string{
			domain.FuelTypeCng,
			domain.FuelTypeDiesel,
			domain.FuelTypeElectric,
			domain.FuelTypePetrol,
			domain.FuelTypeHybrid,
		}
		VehicleTypes := []string{
			domain.CarCategoryMini,
			domain.CarCategoryPremium,
		}

		c.HTML(http.StatusOK, "edit_mini.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"fuel_types":        FuelTypes,
			"car_types":         CarTypes,
			"Car":               car,
			"Brands":            brands,
			"vehicle_types":     VehicleTypes,
			"CurrentPath":       c.Request.URL.Path,
		})

	}
}

func GetChoices(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"car_types": []string{
			domain.CarTypeSuv,
			domain.CarTypeCompactSuv,
			domain.CarTypeConvertible,
			domain.CarTypeHatchback,
			domain.CarTypeMpv,
			domain.CarTypeSedan,
			domain.CarTypeSport,
		},
		"fuel_types": []string{
			domain.FuelTypePetrol,
			domain.FuelTypeDiesel,
			domain.FuelTypeElectric,
			domain.FuelTypeHybrid,
		},
	})

}

// WhatsAppClient handles sending WhatsApp messages via an API

// SendMessage sends a WhatsApp message to the specified number

// AddProduct handles adding a new vehicle and notifying customers
func AddProduct(db *gorm.DB, whatsappClient *services.WhatsAppClient) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var vehicle domain.Vehicle
		var brand domain.Brand
		var images []domain.Image

		// Parse and validate form data
		BrandIDStr := c.PostForm("brand")
		brandID, err := strconv.ParseUint(BrandIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
			return
		}

		if err := db.First(&brand, brandID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand not found"})
			return
		}

		vehicle.BrandID = uint(brandID)
		vehicle.Model = c.PostForm("model")
		year, err := strconv.Atoi(c.PostForm("year"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}
		vehicle.Year = year
		vehicle.Color = c.PostForm("color")
		vehicle.Variant = c.PostForm("variant")
		vehicle.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		vehicle.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		vehicle.Transmission = c.PostForm("transmission")
		vehicle.RegNo = c.PostForm("regno")
		vehicle.Vehicle_type = c.PostForm("vehicle_type")
		vehicle.Status = c.PostForm("status")
		vehicle.Price, _ = strconv.Atoi(c.PostForm("price"))
		vehicle.CarType = c.PostForm("car_type")
		fmt.Println("here is the car type xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", vehicle.CarType)
		vehicle.FuelType = c.PostForm("fuel_type")
		vehicle.Engine_size = c.PostForm("engine_size")
		vehicle.Insurance_date = c.PostForm("insurance_date")
		vehicle.Location = c.PostForm("location")

		// Handle banner image only for premium cars
		if vehicle.Vehicle_type != "Mini" {
			bannerImage, err := c.FormFile("bannerimage")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Banner image is required for premium cars"})
				return
			}
			bannerImagePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), bannerImage.Filename))

			if err := c.SaveUploadedFile(bannerImage, bannerImagePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the banner image"})
				return
			}
			vehicle.BannerImage = "/" + strings.ReplaceAll(bannerImagePath, "\\", "/")
		}

		// Handle file uploads for images
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the form"})
			return
		}
		files := form.File["images[]"]

		for _, file := range files {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}
			imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
			images = append(images, domain.Image{Path: imagePath})
		}
		vehicle.Images = images

		// Add vehicle to the database
		if err := db.Create(&vehicle).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add the car"})
			return
		}

		// Find enquiries matching this vehicle model
		var enquiries []domain.Enquiry
		fmt.Println("here is the enquires")
		if err := db.Where("desired_cars LIKE ?", fmt.Sprintf("%%%s%%", vehicle.Brand.Name)).Find(&enquiries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find matching enquiries"})
			return
		}

		// Create channels for concurrent processing
		messageChannel := make(chan string, len(enquiries))
		fmt.Println("here is the message channel", messageChannel)
		errorChannel := make(chan error, len(enquiries))
		defer close(messageChannel)
		defer close(errorChannel)

		// Send WhatsApp messages concurrently
		for _, enquiry := range enquiries {
			go func(enquiry domain.Enquiry) {
				phoneStr := strconv.Itoa(enquiry.Phone)
				message := fmt.Sprintf("Hello %s, the vehicle %s that you enquired about is now available. Please visit our showroom or contact us for more details.", enquiry.CustomerName, vehicle.Model)
				err := whatsappClient.SendMessage(phoneStr, message)
				if err != nil {
					errorChannel <- fmt.Errorf("Failed to send message to %s: %v", phoneStr, err)
				} else {
					messageChannel <- fmt.Sprintf("Successfully sent message to %s", phoneStr)
				}
			}(enquiry)
		}

		// Collect results from the channels
		for i := 0; i < len(enquiries); i++ {
			select {
			case msg := <-messageChannel:
				log.Println(msg)
			case err := <-errorChannel:
				log.Println(err)
			}
		}

		c.Redirect(http.StatusSeeOther, "/admin/product")
	}
}

func ProductPage(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		var (
			brands            []domain.Brand
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)
		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find the brand"})
		}

		carTypes := []string{domain.CarTypeSuv,
			domain.CarTypeCompactSuv,
			domain.CarTypeConvertible,
			domain.CarTypeHatchback,
			domain.CarTypeMpv,
			domain.CarTypeSedan,
			domain.CarTypeSport,
		}
		fuelTypes := []string{
			domain.FuelTypePetrol,
			domain.FuelTypeDiesel,
			domain.FuelTypeHybrid,
			domain.FuelTypeCng,
			domain.FuelTypeElectric,
		}
		VehicleType := []string{
			domain.CarCategoryPremium,
			domain.CarCategoryMini,
		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count mini cars"})
			return
		}

		c.HTML(http.StatusOK, "product.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"Brands":            brands,
			"CarTypes":          carTypes,
			"FuelTypes":         fuelTypes,
			"Vehicle_type":      VehicleType,
			"CurrentPath":       c.Request.URL.Path,
		})
		fmt.Println("here is teh product")
	}

}

func Logout(c *gin.Context) {
	// Deleting the cookie by setting its expiration to a past time
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.SetCookie("token", "", -1, "/", "", false, true)
	fmt.Println("COME ON EVERY BODY")

	// Check the cookie value after attempting to delete it
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Println("Token deleted successfully:", err)
	} else {
		fmt.Println("Failed to delete token, still present:", token)
	}

	// Redirect to login page after logout
	c.Redirect(http.StatusFound, "/admin/login")
}

func Adminlogin(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	c.Header("Expires", "0")

	// Check for the presence of a token
	token, err := c.Cookie("token")
	fmt.Println("here is the token")
	fmt.Println("here is the login token come on", token)
	if err == nil && token != "" {

		// Validate the token
		valid, _ := utils.ValidateToken(token)
		fmt.Println("here is the valid", valid)
		if valid {
			fmt.Print("yes it is valid dudde")
			// If the token is valid, redirect to the dashboard
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
	}
	// If there's no valid token, render the login page
	c.HTML(http.StatusOK, "login.html", nil)
}

func PremiumCars(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var (
			cars              []domain.Vehicle
			brands            []domain.Brand
			page              int
			totalCount        int64
			limit             int
			offset            int
			EnquiredCustomers []domain.Enquiry
			PremiumCarsTotal  int64
			MiniCarsTotal     int64
			TotalCars         int64
		)

		page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "5"))

		offset = (page - 1) * limit

		if err := db.Preload("Brand").Order("created_at desc").Limit(limit).Offset(offset).Where("vehicle_type = ?", "Premium").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return

		}
		if err := db.Find(&EnquiredCustomers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the enquired customers"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type=?", "Mini").Count(&MiniCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the mini cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type=?", "Premium").Count(&PremiumCarsTotal).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the premium cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Count(&TotalCars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the total cars"})
			return
		}

		if err := db.Model(&domain.Vehicle{}).Where("vehicle_type = ?", "Premium").Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the total count of premium cars"})
			return
		}

		for i := range cars {
			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {

				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the images"})
				return
			}
		}

		noCars := len(cars) == 0

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
		pages := make([]int, totalPages)

		for i := range pages {
			pages[i] = i + 1
		}

		c.HTML(200, "PremiumCars.html", gin.H{
			"enquiredcustomers": len(EnquiredCustomers),
			"premiumcarstotal":  PremiumCarsTotal,
			"minicarstotal":     MiniCarsTotal,
			"totalcars":         TotalCars,
			"Cars":              cars,
			"Brands":            brands,
			"Page":              page,
			"Pages":             pages,
			"totalPages":        totalPages,
			"Limit":             limit,
			"totalCount":        totalCount,
			"NoCars":            noCars,
			"CurrentPath":       c.Request.URL.Path,
		})
	}

}

func AdminLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		var input domain.User
		input.Email = c.PostForm("email")

		fmt.Println("here is the email", input.Email)
		input.Password = c.PostForm("password")

		var user domain.User

		if err := db.Where("email=?", input.Email).First(&user).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
			return
		}

		if !user.IsAdmin {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "user is not admin"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
			return
		}
		token, err := utils.GenerateToken(user.Email)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "token based error"})
			return
		}

		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		fmt.Println("here is the cookie domain", cookieDomain)
		c.SetCookie("token", token, 3600, "/", cookieDomain, false, true)
		fmt.Println("here is the token", token)
		fmt.Print(c.Cookie("token"))

		c.Redirect(http.StatusFound, "/admin/dashboard")

	}
}
