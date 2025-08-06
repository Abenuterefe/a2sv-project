package controllers

import (
	"strconv"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogHandler struct {
	UseCase interfaces.BlogUseCaseInterface
}

func NewBlogHandler(uc interfaces.BlogUseCaseInterface) *BlogHandler {
	return &BlogHandler{UseCase: uc}
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var blog entities.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get authenticated user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.UseCase.CreateBlog(c.Request.Context(), &blog, userID.(string)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, blog)
}

func (h *BlogHandler) GetBlogsByUser(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// Allow optional query parameter to get another user's blogs (for admin or public view)
	targetUserID := c.Query("user_id")
	if targetUserID == "" {
		targetUserID = userID.(string) // Default to authenticated user's blogs
	}

	// parse pagination query parameters
	pageQuery := c.DefaultQuery("page", "1")
	// default and cap: max 5 blogs per page
	limitQuery := c.DefaultQuery("limit", "5")
	page, err := strconv.ParseInt(pageQuery, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQuery, 10, 64)
	if err != nil || limit < 1 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}
	blogs, err := h.UseCase.GetBlogsByUserID(c.Request.Context(), targetUserID, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogs)
}

// GetBlogByID handles GET /blogs/:id
func (h *BlogHandler) GetBlogByID(c *gin.Context) {
	id := c.Param("id")
	blog, err := h.UseCase.GetBlogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Blog not found"})
		return
	}
	c.JSON(200, blog)
}

// UpdateBlog handles PUT /blogs/:id
func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	// First, get the existing blog to preserve all its data
	existingBlog, err := h.UseCase.GetBlogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Blog not found"})
		return
	}

	// Bind the JSON request to the existing blog (this only updates provided fields)
	if err := c.ShouldBindJSON(existingBlog); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure the ID is preserved (shouldn't change during update)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid blog ID"})
		return
	}
	existingBlog.ID = objectID

	// Update the blog
	if err := h.UseCase.UpdateBlog(c.Request.Context(), existingBlog); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, existingBlog)
}

// DeleteBlog handles DELETE /blogs/:id
func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	if err := h.UseCase.DeleteBlog(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
