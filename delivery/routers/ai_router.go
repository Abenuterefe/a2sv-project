package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/ai"
	"github.com/Abenuterefe/a2sv-project/infrastructure/auth"
	"github.com/Abenuterefe/a2sv-project/infrastructure/middlewares"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"github.com/gin-gonic/gin"
)

func AiRoutes(r *gin.Engine) {
	aiService := ai.NewOpenAIService()
	aiUseCase := usecase.NewAIGenerationUseCase(aiService)
	aiController := controllers.NewAIController(aiUseCase)
	jwtService := auth.NewJWTService()
	aiGroup := r.Group("/ai")
	aiGroup.Use(middlewares.AuthMiddleware(jwtService))

	aiGroup.POST("/suggestion", aiController.GenerateBlog)
}