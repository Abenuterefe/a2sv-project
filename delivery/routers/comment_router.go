package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/auth"
	"github.com/Abenuterefe/a2sv-project/infrastructure/middlewares"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommentRoutes initializes the comment-related routes with authentication and authorization.
func CommentRoutes(r *gin.Engine, client *mongo.Client) {
	// Get the comment collection
	commentCollection := client.Database("g6_starter_projectDb").Collection("comments")

	// Initialize JWT service for authentication
	jwtService := auth.NewJWTService()

	// initialization of repo, usecase, and handler
	commentRepo := repository.NewCommentRepositoryMongo(commentCollection)
	commentUseCase := usecase.NewCommentUseCase(commentRepo)
	commentHandler := controllers.NewCommentHandler(commentUseCase)

	// Group routes under /api/v1
	api := r.Group("/api/v1")

	// Public routes (no authentication required)
	api.GET("/blogs/:id/comments", commentHandler.GetCommentsByBlog) // Anyone can view comments on a blog
	api.GET("/comments/:id", commentHandler.GetCommentByID)         // Anyone can view a specific comment

	// Protected routes (authentication required)
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(jwtService))

	// Routes that require authentication
	protected.POST("/blogs/:id/comments", commentHandler.CreateComment) // Create comment (authenticated users only)
	protected.PUT("/comments/:id", commentHandler.UpdateComment)           // Update comment (owner only - checked in handler)
	protected.DELETE("/comments/:id", commentHandler.DeleteComment)        // Delete comment (owner only - checked in handler)
}
