package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/database"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the dependencies and sets up the routes.
func BlogRoutes() *gin.Engine {
	r := gin.Default()

	// Connect to MongoDB
	client, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Get the blog collection
	blogCollection := client.Database("blogdb").Collection("blogs")

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

	return r
}
