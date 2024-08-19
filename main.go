package main

import (
	"encoding/json"
	"fmt"
	"ginapp/config"
	"ginapp/database"
	"ginapp/routes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Custom functions for templates
func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

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

// Convert data to JSON for use in templates
func toJSON(v interface{}) template.HTML {
	bytes, err := json.Marshal(v)
	if err != nil {
		return template.HTML(`{}`)
	}
	return template.HTML(bytes)
}

func main() {
	// Define template functions
	funcMap := template.FuncMap{
		"add":       add,
		"sub":       sub,
		"hasPrefix": strings.HasPrefix,
		"toJson":    toJSON,
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading the config: %v", err)
	}

	// Connect to the database
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	fmt.Println("Database connection successful!")

	// Set up Gin router
	router := gin.Default()
	router.SetFuncMap(funcMap)
	router.Use(CORSMiddleware())
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
