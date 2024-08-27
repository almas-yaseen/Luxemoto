package middleware

import (
	"fmt"
	"ginapp/domain"
	"ginapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func AdminAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("token")
		fmt.Println("here is the token string", tokenString)
		if err != nil {
			fmt.Println("Error retrieving cookie:", err)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		if tokenString == "" {
			fmt.Println("Token is empty")
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			fmt.Println("Error parsing token or token is invalid:", err)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Invalid token claims")
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		var user domain.User
		if err := db.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
			fmt.Println("User not found:", err)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
