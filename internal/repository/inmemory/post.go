package inmemory

import (
	"sync"

	"github.com/tmozzze/SasPosts/internal/domain"
)

type NewInMemoryPostRepository struct {
	mu    sync.RWMutex
	posts map[string]*domain.Post
}

func NewInMemoryPostRepository() *NewInMemoryPostRepository {
	return &inMemoryPostRepository{
		posts: make(map[string]*domain.Post),
	}
}

func (r *InMemoryPostRepository) GetByID(id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, domain.ErrPostNotFound
	}
	return post, nil
}
