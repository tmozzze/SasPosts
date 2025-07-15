package inmemory

import (
	"context"
	"sort"
	"sync"
	"unicode/utf8"

	"github.com/tmozzze/SasPosts/internal/domain"
)

type InMemoryCommentRepository struct {
	mu       sync.RWMutex
	comments map[string]*domain.Comment
}

func NewInMemoryCommentRepository() *InMemoryCommentRepository {
	return &InMemoryCommentRepository{
		comments: make(map[string]*domain.Comment),
	}
}

func (r *InMemoryCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	if utf8.RuneCountInString(comment.Content) > domain.MaxCommentLength {
		return domain.ErrCommentTooLong
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if comment.ParentID == nil {
		comment.Path = comment.ID
		comment.Depth = 0
	} else {
		parent, exists := r.comments[*comment.ParentID]
		if !exists {
			return domain.ErrParentCommentNotFound
		}
		comment.Path = parent.Path + "." + comment.ID
		comment.Depth = parent.Depth + 1
	}

	r.comments[comment.ID] = comment
	return nil
}

func (r *InMemoryCommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, domain.ErrCommentNotFound
	}
	return comment, nil
}

func (r *InMemoryCommentRepository) GetByPost(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID && comment.ParentID == nil {
			results = append(results, comment)
		}
	}

	sortCommentsByCreatedAt(results)

	start := offset
	if start >= len(results) {
		return []*domain.Comment{}, nil
	}
	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

func (r *InMemoryCommentRepository) GetChildren(ctx context.Context, parentID string, limit, offset int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Comment

	for _, comment := range r.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			results = append(results, comment)
		}
	}
	sortCommentsByCreatedAt(results)

	start := offset
	if start >= len(results) {
		return []*domain.Comment{}, nil
	}
	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

func sortCommentsByCreatedAt(comments []*domain.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
}
