package main

import (
	"fmt"
	"ginapp/config"
	"ginapp/database"
	"ginapp/routes"
	"text/template"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func add(x, y int) int {
	return x + y
}

// Sub function to decrement index
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

func main() {
	// Add function to increment index
	password := []byte("jacky")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("jacky"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Hashed Password:", string(hashedPassword))

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err) // nil means it is a match

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading the config", err)
	}

	db, err := database.ConnectDatabase(cfg)
	fmt.Println("come on db", db)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	log.Println("Database connection successf ul!")
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"add": add,
		"sub": sub,
	})
	router.Use(CORSMiddleware())
	router.LoadHTMLGlob("templates/*.html")

	// tmpl := template.Must(template.New("").ParseFiles(files...))

	// router.SetHTMLTemplate(tmpl)

	//admin routes
	adminGroup := router.Group("/admin")
	routes.AdminRoutes(adminGroup, db)

	//user routes

	userGroup := router.Group("/myapp")
	routes.UserRoutes(userGroup, db)

	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
