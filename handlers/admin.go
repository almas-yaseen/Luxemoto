package handlers

import (
	"fmt"
	"ginapp/domain"
	"ginapp/utils"
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

func PremiumCarsDelete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

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
		var brands []domain.Brand
		fmt.Println("here is the brands", brands)

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the brands"})

			return
		}

		c.HTML(http.StatusOK, "brand.html", gin.H{"brands": brands})

	}
}

func EditCarMini(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var (
			cars   []domain.Vehicle
			brands []domain.Brand
		)

		if err := db.Preload("Brand").Where("vehicle_type=?", "Mini").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return
		}

		for i := range cars {

			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the images"})
				return

			}
			c.HTML(200, "MiniCars.html", gin.H{
				"Cars":   cars,
				"Brands": brands,
			})

		}

	}
}

func EditCar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		id := c.Param("id") // fetching the id from the url
		var car domain.Vehicle
		var brands []domain.Brand

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the id"})
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the car"})
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
			"fuel_types":    FuelTypes,
			"car_types":     CarTypes,
			"Car":           car,
			"Brands":        brands,
			"vehicle_types": VehicleTypes,
		})

	}
}
func EditPageMini(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") // fetching the id from the url
		var car domain.Vehicle
		var brands []domain.Brand

		if err := db.Preload("Images").Preload("Brand").First(&car, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the id"})
		}

		if err := db.Find(&brands).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the car"})
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
			"fuel_types":    FuelTypes,
			"car_types":     CarTypes,
			"Car":           car,
			"Brands":        brands,
			"vehicle_types": VehicleTypes,
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
		vehicle.Vehicle_type = c.PostForm("vehicle_type")
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
		VehicleType := []string{
			domain.CarCategoryPremium,
			domain.CarCategoryMini,
		}

		c.HTML(http.StatusOK, "product.html", gin.H{
			"Brands":       brands,
			"CarTypes":     carTypes,
			"FuelTypes":    fuelTypes,
			"Vehicle_type": VehicleType,
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

func PremiumCars(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			cars   []domain.Vehicle
			brands []domain.Brand
		)

		if err := db.Preload("Brand").Where("vehicle_type = ?", "Premium").Find(&cars).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the cars"})
			return

		}

		for i := range cars {
			if err := db.Model(&cars[i]).Association("Images").Find(&cars[i].Images); err != nil {

				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the images"})
				return
			}
		}

		c.HTML(200, "PremiumCars.html", gin.H{
			"Cars":   cars,
			"Brands": brands,
		})
	}

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
