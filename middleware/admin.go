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
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return

		}
		token, err := utils.ParseToken(tokenString)

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		var user domain.User

		if err := db.Where("email=?", claims["email"]).First(&user).Error; err != nil {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
