package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/tmozzze/SasPosts/internal/domain"
	"github.com/tmozzze/SasPosts/utils"
)

type InMemoryPostRepository struct {
	mu    sync.RWMutex
	posts map[string]*domain.Post
}

func NewInMemoryPostRepository() *InMemoryPostRepository {
	return &InMemoryPostRepository{
		posts: make(map[string]*domain.Post),
	}
}

func (r *InMemoryPostRepository) Create(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if post.ID == "" {
		post.ID = utils.GenerateID()
	}
	post.CreatedAt = time.Now()
	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryPostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, domain.ErrPostNotFound
	}
	return post, nil
}

func (r *InMemoryPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *InMemoryPostRepository) ToggleComments(ctx context.Context, postID string, allow bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, exists := r.posts[postID]
	if !exists {
		return domain.ErrPostNotFound
	}

	post.AllowComments = allow
	return nil
}

func (r *InMemoryPostRepository) Update(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.posts[post.ID]
	if !exists {
		return domain.ErrPostNotFound
	}
	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryPostRepository) Delete(ctx context.Context, postID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.posts[postID]
	if !exists {
		return domain.ErrPostNotFound
	}
	delete(r.posts, postID)
	return nil
}

func (r *InMemoryPostRepository) CheckAllowdComments(ctx context.Context, postID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[postID]
	if !exists {
		return false, domain.ErrPostNotFound
	}
	return post.AllowComments, nil
}
