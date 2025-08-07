package controllers

import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	UseCase interfaces.CommentUseCaseInterface
}

func NewCommentHandler(uc interfaces.CommentUseCaseInterface) *CommentHandler {
	return &CommentHandler{UseCase: uc}
}

// CreateComment handles POST /blogs/:blogId/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment entities.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get authenticated user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// Get blog ID from URL parameter
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(400, gin.H{"error": "Blog ID is required"})
		return
	}

	if err := h.UseCase.CreateComment(c.Request.Context(), &comment, userID.(string), blogID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, comment)
}

// GetCommentsByBlog handles GET /blogs/:id/comments
func (h *CommentHandler) GetCommentsByBlog(c *gin.Context) {
	blogID := c.Param("id") // Changed from "blogId" to "id" to match the route parameter
	if blogID == "" {
		c.JSON(400, gin.H{"error": "Blog ID is required"})
		return
	}

	comments, err := h.UseCase.GetCommentsByBlogID(c.Request.Context(), blogID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, comments)
}

// GetCommentByID handles GET /comments/:id
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	id := c.Param("id")
	comment, err := h.UseCase.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}
	c.JSON(200, comment)
}

// UpdateComment handles PUT /comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id := c.Param("id")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// First, get the existing comment to check ownership and preserve data
	existingComment, err := h.UseCase.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership - user can only update their own comments
	if existingComment.UserID != userID.(string) {
		c.JSON(403, gin.H{"error": "You can only modify your own comments"})
		return
	}

	// Bind the JSON request to the existing comment (this only updates provided fields)
	if err := c.ShouldBindJSON(existingComment); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure the ID is preserved (shouldn't change during update)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid comment ID"})
		return
	}
	existingComment.ID = objectID

	// Update the comment
	if err := h.UseCase.UpdateComment(c.Request.Context(), existingComment); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, existingComment)
}

// DeleteComment handles DELETE /comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// First, get the existing comment to check ownership
	existingComment, err := h.UseCase.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership - user can only delete their own comments
	if existingComment.UserID != userID.(string) {
		c.JSON(403, gin.H{"error": "You can only delete your own comments"})
		return
	}

	// Delete the comment
	if err := h.UseCase.DeleteComment(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
