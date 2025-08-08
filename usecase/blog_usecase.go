package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// blogUseCase implements the BlogUseCaseInterface
type blogUseCase struct {
	repo        interfaces.BlogRepositoryInterface
	commentRepo interfaces.CommentRepositoryInterface
}

func NewBlogUseCase(repo interfaces.BlogRepositoryInterface, commentRepo interfaces.CommentRepositoryInterface) interfaces.BlogUseCaseInterface {
	return &blogUseCase{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (u *blogUseCase) CreateBlog(ctx context.Context, blog *entities.Blog, userID string) error {
	// Generate a new ObjectID for the blog
	blog.ID = primitive.NewObjectID()

	// Set the user ID
	blog.UserID = userID

	// Set timestamps
	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now

	return u.repo.CreateBlog(ctx, blog)
}

// GetBlogsByUserID returns paginated blogs for a user
func (u *blogUseCase) GetBlogsByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*entities.Blog, error) {
	return u.repo.GetBlogsByUserID(ctx, userID, page, limit)
}

// GetBlogByID returns a single blog by ID
func (u *blogUseCase) GetBlogByID(ctx context.Context, id string) (*entities.Blog, error) {
	return u.repo.GetBlogByID(ctx, id)
}

// UpdateBlog updates an existing blog
func (u *blogUseCase) UpdateBlog(ctx context.Context, blog *entities.Blog) error {
	// Update timestamp
	blog.UpdatedAt = time.Now()
	return u.repo.UpdateBlog(ctx, blog)
}

// DeleteBlog removes a blog by ID
func (u *blogUseCase) DeleteBlog(ctx context.Context, id string) error {
	return u.repo.DeleteBlog(ctx, id)
}

// calculatePopularityScore calculates the popularity score for a blog
func (u *blogUseCase) calculatePopularityScore(blog *entities.Blog) float64 {
	// Get comment count
	commentCount, _ := u.commentRepo.GetCommentCountByBlogID(context.Background(), blog.ID.Hex())

	// Base popularity score using your algorithm
	score := float64(blog.LikeCount*3) + float64(commentCount*5) + float64(blog.ViewCount)*0.1 + float64(blog.DislikeCount*-2)

	// Recency boost
	recencyBoost := u.calculateRecencyBoost(blog.CreatedAt)

	return score + recencyBoost
}

// calculateRecencyBoost calculates bonus points for recent posts
func (u *blogUseCase) calculateRecencyBoost(createdAt time.Time) float64 {
	daysSinceCreation := time.Since(createdAt).Hours() / 24

	if daysSinceCreation <= 1 {
		return 50 // Very recent posts
	} else if daysSinceCreation <= 7 {
		return 20 // Posts from this week
	} else if daysSinceCreation <= 30 {
		return 5 // Posts from this month
	}
	return 0 // Older posts
}

// GetPopularBlogs retrieves blogs sorted by popularity score
func (u *blogUseCase) GetPopularBlogs(ctx context.Context, limit int64) ([]*entities.BlogWithPopularity, error) {
	blogs, err := u.repo.GetAllBlogs(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to BlogWithPopularity and calculate scores
	popularBlogs := make([]*entities.BlogWithPopularity, 0, len(blogs))
	for _, blog := range blogs {
		commentCount, _ := u.commentRepo.GetCommentCountByBlogID(ctx, blog.ID.Hex())

		popularBlog := &entities.BlogWithPopularity{
			ID:              blog.ID,
			Title:           blog.Title,
			Content:         blog.Content,
			UserID:          blog.UserID,
			LikeCount:       blog.LikeCount,
			DislikeCount:    blog.DislikeCount,
			ViewCount:       blog.ViewCount,
			CommentCount:    int(commentCount),
			PopularityScore: u.calculatePopularityScore(blog),
			CreatedAt:       blog.CreatedAt,
			UpdatedAt:       blog.UpdatedAt,
		}
		popularBlogs = append(popularBlogs, popularBlog)
	}

	// Sort by popularity score (descending)
	sort.Slice(popularBlogs, func(i, j int) bool {
		return popularBlogs[i].PopularityScore > popularBlogs[j].PopularityScore
	})

	// Apply limit
	if limit > 0 && int64(len(popularBlogs)) > limit {
		popularBlogs = popularBlogs[:limit]
	}

	return popularBlogs, nil
}

// FilterBlogs filters blogs based on provided criteria
func (u *blogUseCase) FilterBlogs(ctx context.Context, filter *entities.BlogFilter) (*entities.FilterResponse, error) {
	// Validate filter criteria
	if filter.DateFrom != nil && filter.DateTo != nil {
		if filter.DateFrom.After(*filter.DateTo) {
			return nil, errors.New("date_from cannot be after date_to")
		}
	}

	// Set default values
	if filter.Limit == 0 {
		filter.Limit = 20 // default limit
	}

	// Validate popularity sort options
	if filter.PopularitySort != "" {
		validSortTypes := map[string]bool{
			"views":      true,
			"likes":      true,
			"dislikes":   true,
			"engagement": true,
		}
		if !validSortTypes[filter.PopularitySort] {
			return nil, errors.New("invalid popularity_sort value. Valid values: views, likes, dislikes, engagement")
		}
	}

	// Validate sort order
	if filter.SortOrder != "" && filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		return nil, errors.New("invalid sort_order value. Valid values: asc, desc")
	}

	// Get filtered blogs from repository
	blogs, totalCount, err := u.repo.FilterBlogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate page info
	page := 1
	if filter.Skip > 0 && filter.Limit > 0 {
		page = (filter.Skip / filter.Limit) + 1
	}

	// Create response
	response := &entities.FilterResponse{
		Blogs:      blogs,
		Count:      len(blogs),
		TotalCount: totalCount,
		Page:       page,
		Limit:      filter.Limit,
	}

	return response, nil
}

// SearchBlogs searches for blogs based on title and/or author name
func (u *blogUseCase) SearchBlogs(ctx context.Context, search *entities.BlogSearch) (*entities.SearchResponse, error) {
	// Validate search criteria - at least one search parameter must be provided
	if search.Title == "" && search.Author == "" {
		return nil, errors.New("at least one search parameter (title or author) must be provided")
	}
	
	// Set default values
	if search.Limit == 0 {
		search.Limit = 20 // default limit
	}
	
	// Validate limit and skip
	if search.Limit < 0 {
		return nil, errors.New("limit must be non-negative")
	}
	if search.Skip < 0 {
		return nil, errors.New("skip must be non-negative")
	}
	
	// Get search results from repository
	blogs, totalCount, err := u.repo.SearchBlogs(ctx, search)
	if err != nil {
		return nil, err
	}
	
	// Create response
	response := &entities.SearchResponse{
		Blogs:      blogs,
		Count:      len(blogs),
		TotalCount: totalCount,
		Query:      search,
	}
	
	return response, nil
}
