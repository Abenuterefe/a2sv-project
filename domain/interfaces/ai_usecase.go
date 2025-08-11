package interfaces
import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
)
type AIGenerationUseCaseInterface interface {
	GenerateBlog(prompt string) (*entities.BlogResponse, error)
}