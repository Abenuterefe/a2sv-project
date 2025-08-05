package controllers

import (
	"a2sv-project/Domain/entities"
	interfaces "a2sv-project/Domain/interfaces"
	"strconv"

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

	// TODO: Get user ID from context (auth middleware)
	// For now, using a default user ID - replace this with actual auth logic
	userID := "default-user-id" // This should come from authentication

	if err := h.UseCase.CreateBlog(c.Request.Context(), &blog, userID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, blog)
}

func (h *BlogHandler) GetBlogsByUser(c *gin.Context) {
	// TODO: Get user ID from context (auth middleware)
	userID := c.Query("user_id")
	if userID == "" {
		// Use default user ID if not provided
		userID = "default-user-id"
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
	blogs, err := h.UseCase.GetBlogsByUserID(c.Request.Context(), userID, page, limit)
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
	var blog entities.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	// Convert string id to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid blog ID"})
		return
	}
	blog.ID = objectID
	if err := h.UseCase.UpdateBlog(c.Request.Context(), &blog); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blog)
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
