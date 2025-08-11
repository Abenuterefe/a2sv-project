package usecase

import (
	"context"
	"testing"
	"time"

	repoMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLikeBlog_ToggleOff(t *testing.T) {
	t.Parallel()
	interRepo := repoMocks.NewBlogInteractionRepositoryInterface(t)
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	uc := NewBlogInteractionUseCase(interRepo, blogRepo)

	// already liked => remove like, decrement like counter
	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "like").Return(true, nil)
	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "dislike").Return(false, nil)
	interRepo.On("RemoveInteraction", mock.Anything, "b1", "u1", "like").Return(nil)
	blogRepo.On("UpdateBlogCounters", mock.Anything, "b1", -1, 0, 0).Return(nil)

	err := uc.LikeBlog(context.Background(), "b1", "u1")
	assert.NoError(t, err)
}

func TestLikeBlog_SwitchFromDislike(t *testing.T) {
	t.Parallel()
	interRepo := repoMocks.NewBlogInteractionRepositoryInterface(t)
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	uc := NewBlogInteractionUseCase(interRepo, blogRepo)

	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "like").Return(false, nil)
	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "dislike").Return(true, nil)
	interRepo.On("RemoveInteraction", mock.Anything, "b1", "u1", "dislike").Return(nil)
	interRepo.On("AddInteraction", mock.Anything, mock.MatchedBy(func(i interface{}) bool { return true })).Return(nil)
	blogRepo.On("UpdateBlogCounters", mock.Anything, "b1", 1, -1, 0).Return(nil)

	err := uc.LikeBlog(context.Background(), "b1", "u1")
	assert.NoError(t, err)
}

func TestDislikeBlog_New(t *testing.T) {
	t.Parallel()
	interRepo := repoMocks.NewBlogInteractionRepositoryInterface(t)
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	uc := NewBlogInteractionUseCase(interRepo, blogRepo)

	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "like").Return(false, nil)
	interRepo.On("HasInteraction", mock.Anything, "b1", "u1", "dislike").Return(false, nil)
	interRepo.On("AddInteraction", mock.Anything, mock.MatchedBy(func(i interface{}) bool { return true })).Return(nil)
	blogRepo.On("UpdateBlogCounters", mock.Anything, "b1", 0, 1, 0).Return(nil)

	err := uc.DislikeBlog(context.Background(), "b1", "u1")
	assert.NoError(t, err)
}

func TestViewBlog_Anonymous_Debounce(t *testing.T) {
	t.Parallel()
	interRepo := repoMocks.NewBlogInteractionRepositoryInterface(t)
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	uc := NewBlogInteractionUseCase(interRepo, blogRepo)

	// First call indicates recent view exists -> no increment, no add
	blogID := "507f1f77bcf86cd799439011" // valid ObjectID hex
	interRepo.On("HasRecentView", mock.Anything, blogID, "anonymous", "1.1.1.1", "agent").Return(true, nil)

	err := uc.ViewBlog(context.Background(), blogID, "", "1.1.1.1", "agent")
	assert.NoError(t, err)
}

func TestViewBlog_Anonymous_AddsAndIncrements(t *testing.T) {
	t.Parallel()
	interRepo := repoMocks.NewBlogInteractionRepositoryInterface(t)
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	uc := NewBlogInteractionUseCase(interRepo, blogRepo)

	blogID := "507f1f77bcf86cd799439011" // valid ObjectID hex
	interRepo.On("HasRecentView", mock.Anything, blogID, "anonymous", "1.1.1.1", "agent").Return(false, nil)
	interRepo.On("AddInteraction", mock.Anything, mock.MatchedBy(func(i interface{}) bool { return true })).Return(nil)
	blogRepo.On("UpdateBlogCounters", mock.Anything, blogID, 0, 0, 1).Return(nil)

	err := uc.ViewBlog(context.Background(), blogID, "", "1.1.1.1", "agent")
	assert.NoError(t, err)
}

// avoid unused import lint by touching time
var _ = time.Now
