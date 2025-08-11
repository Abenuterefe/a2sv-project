package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	repoMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Basic contract for NewBlogUseCase validations without hitting repositories
func TestFilterBlogs_InvalidDateRange(t *testing.T) {
	t.Parallel()

	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)

	uc := NewBlogUseCase(blogRepo, commentRepo)

	// date_from after date_to should be rejected
	df := time.Now().Add(24 * time.Hour)
	dt := time.Now()
	_, err := uc.FilterBlogs(context.Background(), &entities.BlogFilter{
		DateFrom: &df,
		DateTo:   &dt,
	})

	assert.Error(t, err)
	assert.Equal(t, "date_from cannot be after date_to", err.Error())
}

func TestSearchBlogs_MissingParams(t *testing.T) {
	t.Parallel()

	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)

	uc := NewBlogUseCase(blogRepo, commentRepo)

	// both title and author are empty
	resp, err := uc.SearchBlogs(context.Background(), &entities.BlogSearch{})
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "at least one search parameter (title or author) must be provided", err.Error())
}

// A tiny happy-path smoke test for SearchBlogs pagination defaults
func TestSearchBlogs_DefaultLimitAndNonNegative(t *testing.T) {
	t.Parallel()

	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)

	uc := NewBlogUseCase(blogRepo, commentRepo)

	blogRepo.On("SearchBlogs", mock.Anything, mock.MatchedBy(func(s *entities.BlogSearch) bool {
		return s.Title == "Go" && s.Limit == 20 && s.Skip == 0
	})).Return([]*entities.BlogWithAuthor{}, int64(0), nil)

	resp, err := uc.SearchBlogs(context.Background(), &entities.BlogSearch{Title: "Go"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, resp.Count)
	assert.Equal(t, int64(0), resp.TotalCount)
}

func TestFilterBlogs_InvalidPopularitySort(t *testing.T) {
	t.Parallel()

	uc := NewBlogUseCase(repoMocks.NewBlogRepositoryInterface(t), repoMocks.NewCommentRepositoryInterface(t))
	_, err := uc.FilterBlogs(context.Background(), &entities.BlogFilter{PopularitySort: "unknown"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid popularity_sort value")
}

func TestFilterBlogs_InvalidSortOrder(t *testing.T) {
	t.Parallel()

	uc := NewBlogUseCase(repoMocks.NewBlogRepositoryInterface(t), repoMocks.NewCommentRepositoryInterface(t))
	_, err := uc.FilterBlogs(context.Background(), &entities.BlogFilter{SortOrder: "up"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_order value")
}

func TestFilterBlogs_HappyPath_PageAndCount(t *testing.T) {
	t.Parallel()

	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewBlogUseCase(blogRepo, commentRepo)

	blogs := []*entities.Blog{{Title: "A"}, {Title: "B"}}
	blogRepo.On("FilterBlogs", mock.Anything, mock.MatchedBy(func(f *entities.BlogFilter) bool {
		return f.Limit == 10 && f.Skip == 10 // page=2, limit=10 => skip=10
	})).Return(blogs, int64(25), nil)

	// page -> skip conversion logic executed in handler, not in usecase; here we pass Skip directly
	resp, err := uc.FilterBlogs(context.Background(), &entities.BlogFilter{Limit: 10, Skip: 10})
	assert.NoError(t, err)
	assert.Equal(t, 2, resp.Count)
	assert.Equal(t, int64(25), resp.TotalCount)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 2, resp.Page)
}

func TestSearchBlogs_NegativeLimitSkip(t *testing.T) {
	t.Parallel()

	uc := NewBlogUseCase(repoMocks.NewBlogRepositoryInterface(t), repoMocks.NewCommentRepositoryInterface(t))

	_, err := uc.SearchBlogs(context.Background(), &entities.BlogSearch{Title: "x", Limit: -1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be non-negative")

	_, err = uc.SearchBlogs(context.Background(), &entities.BlogSearch{Title: "x", Skip: -1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skip must be non-negative")
}

func TestGetPopularBlogs_SortsAndLimits(t *testing.T) {
	t.Parallel()

	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewBlogUseCase(blogRepo, commentRepo)

	// Create 3 blogs with different metrics
	b1 := &entities.Blog{ID: primitive.NewObjectID(), Title: "Old but many views", ViewCount: 1000, LikeCount: 10, DislikeCount: 1, CreatedAt: time.Now().Add(-40 * 24 * time.Hour)}
	b2 := &entities.Blog{ID: primitive.NewObjectID(), Title: "New and liked", ViewCount: 50, LikeCount: 100, DislikeCount: 0, CreatedAt: time.Now().Add(-12 * time.Hour)}
	b3 := &entities.Blog{ID: primitive.NewObjectID(), Title: "Average", ViewCount: 200, LikeCount: 20, DislikeCount: 2, CreatedAt: time.Now().Add(-10 * 24 * time.Hour)}

	blogRepo.On("GetAllBlogs", mock.Anything).Return([]*entities.Blog{b1, b2, b3}, nil)

	// Popularity uses comment counts twice per blog (score + explicit field)
	// We'll return comment counts per blog consistently
	counts := map[string]int64{b1.ID.Hex(): 2, b2.ID.Hex(): 30, b3.ID.Hex(): 5}
	commentRepo.On("GetCommentCountByBlogID", mock.Anything, b1.ID.Hex()).Return(counts[b1.ID.Hex()], nil).Twice()
	commentRepo.On("GetCommentCountByBlogID", mock.Anything, b2.ID.Hex()).Return(counts[b2.ID.Hex()], nil).Twice()
	commentRepo.On("GetCommentCountByBlogID", mock.Anything, b3.ID.Hex()).Return(counts[b3.ID.Hex()], nil).Twice()

	popular, err := uc.GetPopularBlogs(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, popular, 2)
	// Expect the very recent & highly liked b2 to be first
	assert.Equal(t, "New and liked", popular[0].Title)
}

func TestCreateBlog_SetsFieldsAndCallsRepo(t *testing.T) {
	t.Parallel()
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewBlogUseCase(blogRepo, commentRepo)

	// Expect CreateBlog with blog having ID, userID and timestamps set
	blogRepo.On("CreateBlog", mock.Anything, mock.MatchedBy(func(b *entities.Blog) bool {
		return b.ID.Hex() != "" && b.UserID == "u1" && !b.CreatedAt.IsZero() && !b.UpdatedAt.IsZero()
	})).Return(nil)

	err := uc.CreateBlog(context.Background(), &entities.Blog{Title: "t"}, "u1")
	assert.NoError(t, err)
}

func TestUpdateBlog_UpdatesTimestampAndCallsRepo(t *testing.T) {
	t.Parallel()
	blogRepo := repoMocks.NewBlogRepositoryInterface(t)
	commentRepo := repoMocks.NewCommentRepositoryInterface(t)
	uc := NewBlogUseCase(blogRepo, commentRepo)

	before := time.Now().Add(-time.Minute)
	blog := &entities.Blog{Title: "t", UpdatedAt: before}

	blogRepo.On("UpdateBlog", mock.Anything, mock.MatchedBy(func(b *entities.Blog) bool {
		return b.UpdatedAt.After(before) || b.UpdatedAt.Equal(before) == false
	})).Return(nil)

	err := uc.UpdateBlog(context.Background(), blog)
	assert.NoError(t, err)
}
