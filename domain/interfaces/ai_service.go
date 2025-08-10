package interfaces

type AIGenerationInterface interface {
	GenerateBlog(prompt string) (string, error)
}