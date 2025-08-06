package usecase

import (
	"context"
	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
)

// BlogUseCaseInterface is now defined in Domain/interfaces/blog_usecase_interface.go
type blogUseCase struct {
	repo repository.BlogRepository
}

func NewBlogUseCase(repo repository.BlogRepository) interfaces.BlogUseCaseInterface {
	   return &blogUseCase{repo: repo}
}

func (u *blogUseCase) CreateBlog(ctx context.Context, blog *entities.Blog) error {
	   return u.repo.CreateBlog(ctx, blog)
}

func (u *blogUseCase) GetBlogsByUser(ctx context.Context, userID string) ([]*entities.Blog, error) {
	   return u.repo.GetBlogsByAuthorID(ctx, userID)
}
