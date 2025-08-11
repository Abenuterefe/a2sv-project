package interfaces
import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
)
type AIGenerationInterface interface {
	GenerateBlog(prompt string) (*entities.BlogResponse, error)
}