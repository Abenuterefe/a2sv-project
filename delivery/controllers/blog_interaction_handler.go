package controllers

import (
	"net/http"

	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type BlogInteractionHandler struct {
	UseCase interfaces.BlogInteractionUseCaseInterface
}

func NewBlogInteractionHandler(uc interfaces.BlogInteractionUseCaseInterface) *BlogInteractionHandler {
	return &BlogInteractionHandler{UseCase: uc}
}

// LikeBlog handles POST /blogs/:id/like
func (h *BlogInteractionHandler) LikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if err := h.UseCase.LikeBlog(c.Request.Context(), blogID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}

// DislikeBlog handles POST /blogs/:id/dislike
func (h *BlogInteractionHandler) DislikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if err := h.UseCase.DislikeBlog(c.Request.Context(), blogID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog disliked successfully"})
}

// ViewBlog handles POST /blogs/:id/view
func (h *BlogInteractionHandler) ViewBlog(c *gin.Context) {
	blogID := c.Param("id")
	
	// Get user ID (empty string if not authenticated)
	userID, exists := c.Get("userID")
	var uid string
	if exists && userID != nil {
		uid = userID.(string)
	}
	
	// Get IP address and User-Agent for anonymous tracking
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	if err := h.UseCase.ViewBlog(c.Request.Context(), blogID, uid, ipAddress, userAgent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Blog view recorded"})
}
