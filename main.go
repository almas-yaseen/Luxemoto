package main

import (
	"fmt"
	"ginapp/config"
	"ginapp/database"
	"ginapp/routes"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func main() {
	password := "luxemoto@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing password", err)
		return
	}
	fmt.Println("hashed password", string(hashedPassword))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading the config", err)
	}

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	log.Println("Database connection successful!")
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.LoadHTMLGlob("templates/*")

	// tmpl := template.Must(template.New("").ParseFiles(files...))

	// router.SetHTMLTemplate(tmpl)
	adminGroup := router.Group("/admin")

	routes.AdminRoutes(adminGroup, db)

	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
