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

func BlogInteractionRoutes(r *gin.Engine, client *mongo.Client) {
	interactionCollection := client.Database("g6_starter_projectDb").Collection("blog_interactions")
	blogCollection := client.Database("g6_starter_projectDb").Collection("blogs")
	jwtService := auth.NewJWTService()

	interactionRepo := repository.NewBlogInteractionRepositoryMongo(interactionCollection)
	blogRepo := repository.NewBlogRepositoryMongo(blogCollection)

	interactionUseCase := usecase.NewBlogInteractionUseCase(interactionRepo, blogRepo)
	interactionHandler := controllers.NewBlogInteractionHandler(interactionUseCase)

	api := r.Group("/api/v1")

	// Like/dislike require authentication
	protected := api.Group("/blogs")
	protected.Use(middlewares.AuthMiddleware(jwtService))
	protected.POST(":id/like", interactionHandler.LikeBlog)
	protected.POST(":id/dislike", interactionHandler.DislikeBlog)

	// Views can be anonymous (no auth required)
	api.POST("/blogs/:id/view", interactionHandler.ViewBlog)
}
