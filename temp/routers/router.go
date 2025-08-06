package routers

import (
	"a2sv-project/Delivery/controllers"
	"a2sv-project/Infrastructure/database"
	repository "a2sv-project/Repository"
	usecase "a2sv-project/Usecase"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the dependencies and sets up the routes.
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// initialization of repo, usecase, and handler
	blogRepo := repository.NewBlogRepositoryMongo(database.GetCollection("blogs"))
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
