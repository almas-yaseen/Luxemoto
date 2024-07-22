package handlers

import (
	"fmt"
	"ginapp/domain"
	"ginapp/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

func AddProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vehicle domain.Vehicle
		var brand domain.Brand
		var images []domain.Image

		BrandIDStr := c.PostForm("brand")
		fmt.Println("here is the brand str", BrandIDStr)
		brandID, err := strconv.ParseUint(BrandIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand id"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		vehicle.Year = year
		vehicle.Color = c.PostForm("color")
		vehicle.Variant = c.PostForm("variant")
		vehicle.Kms, _ = strconv.Atoi(c.PostForm("kms"))
		vehicle.Ownership, _ = strconv.Atoi(c.PostForm("ownership"))
		vehicle.Transmission = c.PostForm("transmission")
		vehicle.RegNo = c.PostForm("regno")
		vehicle.Status = c.PostForm("status")
		vehicle.Price, _ = strconv.Atoi(c.PostForm("price"))
		vehicle.CarType = c.PostForm("car_type")
		vehicle.FuelType = c.PostForm("fuel_type")
		vehicle.Engine_size = c.PostForm("engine_size")
		vehicle.Insurance_date = c.PostForm("insurance_date")
		vehicle.Location = c.PostForm("location")

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get the form"})
			return
		}
		files := form.File["images[]"]
		fmt.Println("here is the files", files)

		for _, file := range files {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			uploadPath := filepath.Join("uploads", filename)
			if err := c.SaveUploadedFile(file, uploadPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save an image"})
				return
			}
			imagePath := "/" + strings.ReplaceAll(uploadPath, "\\", "/")
			images = append(images, domain.Image{Path: imagePath})
			vehicle.Images = images

			if err := db.Create(&vehicle).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the car"})
				return
			}

			if err := db.Preload("Brand").First(&vehicle, vehicle.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load the car brand"})
				return
			}

			c.Redirect(http.StatusSeeOther, "/admin/dashboard")

		}

	}
}

func ProductPage(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var brands []domain.Brand
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

		c.HTML(http.StatusOK, "product.html", gin.H{
			"Brands":    brands,
			"CarTypes":  carTypes,
			"FuelTypes": fuelTypes,
		})
		fmt.Println("here is teh product")
	}

}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}

func Adminlogin(c *gin.Context) {
	// Check for the presence of a token
	token, err := c.Cookie("token")
	fmt.Println("here is the login token", token)
	if err == nil && token != "" {
		// Validate the token
		valid, _ := utils.ValidateToken(token)
		fmt.Println("here is the valid", valid)
		if valid {
			// If the token is valid, redirect to the dashboard
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
	}
	// If there's no valid token, render the login page
	c.HTML(http.StatusOK, "login.html", nil)
}
func AdminDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{})
}

func PremiumCars(c *gin.Context) {
	c.HTML(200, "PremiumCars.html", gin.H{})
}

func AdminLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		fmt.Println("here is the token", token)
		c.Redirect(http.StatusFound, "/admin/dashboard")

	}
}
