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

// BlogRoutes initializes the blog-related routes with authentication and authorization.
func BlogRoutes(r *gin.Engine, client *mongo.Client) {
	// Get collections
	blogCollection := client.Database("g6_starter_projectDb").Collection("blogs")
	commentCollection := client.Database("g6_starter_projectDb").Collection("comments")

	// Initialize JWT service for authentication
	jwtService := auth.NewJWTService()

	// initialization of repositories, usecase, and handler
	blogRepo := repository.NewBlogRepositoryMongo(blogCollection)
	commentRepo := repository.NewCommentRepositoryMongo(commentCollection)
	blogUseCase := usecase.NewBlogUseCase(blogRepo, commentRepo)
	blogHandler := controllers.NewBlogHandler(blogUseCase)

	// Group routes under /api/v1
	api := r.Group("/api/v1")

	// Public routes (no authentication required)
	api.GET("/blogs/:id", blogHandler.GetBlogByID) // Anyone can view a specific blog
	api.GET("/blogs/popular", blogHandler.GetPopularBlogs) // Anyone can view popular blogs

	// Protected routes (authentication required)
	protected := api.Group("/blogs")
	protected.Use(middlewares.AuthMiddleware(jwtService))

	// Routes that require authentication
	protected.POST("", blogHandler.CreateBlog)    // Create blog (authenticated users only)
	protected.GET("", blogHandler.GetBlogsByUser) // Get user's blogs (authenticated users only)

	// Routes that require authentication + ownership verification
	ownershipProtected := protected.Group("")
	ownershipProtected.Use(middlewares.BlogOwnershipMiddleware(blogUseCase))
	ownershipProtected.PUT("/:id", blogHandler.UpdateBlog)    // Update blog (owner only)
	ownershipProtected.DELETE("/:id", blogHandler.DeleteBlog) // Delete blog (owner only)
}
