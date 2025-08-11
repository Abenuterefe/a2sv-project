package usecase

import (
	"context"
	"testing"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	repoMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateComment_InvalidBlogID(t *testing.T) {
	t.Parallel()
	repo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewCommentUseCase(repo)

	err := uc.CreateComment(context.Background(), &entities.Comment{Content: "hi"}, "u1", "badid")
	assert.Error(t, err)
}

func TestCreateComment_Success(t *testing.T) {
	t.Parallel()
	repo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewCommentUseCase(repo)

	repo.On("CreateComment", mock.Anything, mock.Anything).Return(nil)

	err := uc.CreateComment(context.Background(), &entities.Comment{Content: "hi"}, "u1", "507f1f77bcf86cd799439011")
	assert.NoError(t, err)
}
