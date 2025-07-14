// graph/resolver.go

package graph

import (
	"sync"

	"github.com/tmozzze/SasPosts/internal/domain"
	"github.com/tmozzze/SasPosts/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostRepo         repository.PostRepository
	CommentRepo      repository.CommentRepository
	commentObservers map[string][]chan *domain.Comment
	mu               sync.Mutex
}

func NewResolver(postRepo repository.PostRepository, commentRepo repository.CommentRepository) *Resolver {
	return &Resolver{
		PostRepo:         postRepo,
		CommentRepo:      commentRepo,
		commentObservers: make(map[string][]chan *domain.Comment),
	}
}
