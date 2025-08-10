package controllers

import (
	"net/http"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiUseCase interfaces.AIGenerationUseCaseInterface
}

func NewAIController(aiUseCase interfaces.AIGenerationUseCaseInterface) *AIController {
	return &AIController{aiUseCase}
}

func (ctrl *AIController) GenerateBlog(c *gin.Context) {
	var req entities.PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := ctrl.aiUseCase.GenerateBlog(req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": result})
}
