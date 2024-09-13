// main.go
package main

import (
	"encoding/json"
	"fmt"
	"ginapp/config"
	"ginapp/database"
	"ginapp/middleware"
	"ginapp/routes"
	"html/template"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Custom template functions
func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func toJSON(v interface{}) template.HTML {
	bytes, err := json.Marshal(v)
	if err != nil {
		return template.HTML("{}")
	}
	return template.HTML(bytes)
}

// Load configuration
func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {

	password := "almas"
	hashsedPassword, err := HashPassword(password)

	if err != nil {
		log.Fatalf("Failed hash password %v", err)
	}

	fmt.Println("here is the hashed password", hashsedPassword)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}

	// Connect to the database
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Database connection successful!")

	// Initialize WhatsApp client

	// Set up Gin router
	router := gin.Default()

	// Define template functions
	funcMap := template.FuncMap{
		"add":       add,
		"sub":       sub,
		"hasPrefix": strings.HasPrefix,
		"toJson":    toJSON,
	}
	router.SetFuncMap(funcMap)
	router.Use(middleware.CORSMiddleware())
	router.LoadHTMLGlob("templates/*.html")

	// Define routes
	adminGroup := router.Group("/admin")
	routes.AdminRoutes(adminGroup, db)

	userGroup := router.Group("/myapp")
	routes.UserRoutes(userGroup, db)

	// Serve static files
	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	// Start the server
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
