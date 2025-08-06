package middlewares

import (
	"net/http"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/gin-gonic/gin"
)

// Admin middle ware ristrict access to admin only
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != string(entities.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// useronly middleware ristrict access to non admin users
func UserOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != string(entities.RoleUser) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
