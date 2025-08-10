package interfaces

type AIGenerationUseCaseInterface interface {
	GenerateBlog(prompt string) (string, error)
}