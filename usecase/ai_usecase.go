package usecase

import "github.com/Abenuterefe/a2sv-project/domain/interfaces"

type AIGenerationUseCase struct {
	aiService interfaces.AIGenerationInterface
}

func NewAIGenerationUseCase(aiService interfaces.AIGenerationInterface) interfaces.AIGenerationUseCaseInterface {
	return &AIGenerationUseCase{aiService}
}

func (u *AIGenerationUseCase) GenerateBlog(prompt string) (string, error) {
	return u.aiService.GenerateBlog(prompt)
}