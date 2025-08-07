package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type BlogInteractionRepositoryInterface interface {
	AddInteraction(ctx context.Context, interaction *entities.BlogInteraction) error
	RemoveInteraction(ctx context.Context, blogID string, userID string, interactionType string) error
	HasInteraction(ctx context.Context, blogID string, userID string, interactionType string) (bool, error)
}
