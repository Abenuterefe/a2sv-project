package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// BlogRoutes initializes the blog-related routes.
func BlogRoutes(r *gin.Engine, client *mongo.Client) {
	// Get the blog collection
	blogCollection := client.Database("g6_starter_projectDb").Collection("blogs")

	// initialization of repo, usecase, and handler
	blogRepo := repository.NewBlogRepositoryMongo(blogCollection)
	blogUseCase := usecase.NewBlogUseCase(blogRepo)
	blogHandler := controllers.NewBlogHandler(blogUseCase)

	// Group routes under /api/v1
	api := r.Group("/api/v1")

	// Blog routes
	api.POST("/blogs", blogHandler.CreateBlog)
	api.GET("/blogs", blogHandler.GetBlogsByUser)
	api.GET("/blogs/:id", blogHandler.GetBlogByID)
	api.PUT("/blogs/:id", blogHandler.UpdateBlog)
	api.DELETE("/blogs/:id", blogHandler.DeleteBlog)
}
