package middleware

import (
	"admin-dashboard-FP/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte(os.Getenv("SecretKey"))
		tokenString := c.GetHeader("Authorization")

		dataToken := &models.DataToken{}

		token, err := jwt.ParseWithClaims(tokenString, dataToken, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("Token", dataToken)
		c.Next()
	}
}

func AuthorAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		dataToken, exist := c.Get("Token")

		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		if u, ok := dataToken.(*models.DataToken); ok && u.Role == "admin" {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbbiden access"})
			c.Abort()
		}
	}
}
