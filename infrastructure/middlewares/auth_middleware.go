package middlewares


import (
	"net/http"
	"strings"

	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService interfaces.AuthService) gin.HandlerFunc {
	return  func(c *gin.Context) {
		// exracact token from authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return 
		}

		// Extract token 
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return 
		}

		tokenString := parts[1]

		// verify token
		claims,err := authService.VerifyToken(tokenString,true)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return 
		}

		// Attach user info to context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}