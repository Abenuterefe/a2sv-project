package middlewares

import (
	"net/http"

	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

// BlogOwnershipMiddleware checks if the authenticated user owns the blog they're trying to modify
func BlogOwnershipMiddleware(blogUseCase interfaces.BlogUseCaseInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get blog ID from URL parameter
		blogID := c.Param("id")
		if blogID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
			c.Abort()
			return
		}

		// Fetch the blog to check ownership
		blog, err := blogUseCase.GetBlogByID(c.Request.Context(), blogID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
			c.Abort()
			return
		}

		// Check if the user owns this blog (debug logging)
		userIDStr := userID.(string)
		if blog.UserID != userIDStr {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You can only modify your own blogs",
				"debug": gin.H{
					"blogUserID": blog.UserID,
					"authUserID": userIDStr,
				},
			})
			c.Abort()
			return
		}

		// User owns the blog, continue
		c.Next()
	}
}
