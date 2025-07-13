package repository

import (
	"context"

	"github.com/tmozzze/SasPosts/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	GetAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, postID string) error
	CheckAllowedComments(ctx context.Context, postID string) (bool, error)
	ToggleComments(ctx context.Context, postID string, allow bool) error
}
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	GetByPost(ctx context.Context, postID string, limit int, offset int) ([]*domain.Comment, error)
	GetChildren(ctx context.Context, parentID string, limit int, offset int) ([]*domain.Comment, error)
	CountByPost(ctx context.Context, postID string) (int, error)
	CountChildren(ctx context.Context, parentID string) (int, error)
}
