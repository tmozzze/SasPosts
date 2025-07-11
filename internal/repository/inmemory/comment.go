package inmemory

import (
	"context"
	"sort"
	"sync"
	"unicode/utf8"

	"github.com/tmozzze/SasPosts/internal/domain"
)

type inMemoryCommentRepository struct {
	mu       sync.RWMutex
	comments map[string]*domain.Comment
	byPost   map[string][]string
	byParent map[string][]string
}

func NewInMemoryCommentRepository() *inMemoryCommentRepository {
	return &inMemoryCommentRepository{
		comments: make(map[string]*domain.Comment),
		byPost:   make(map[string][]string),
		byParent: make(map[string][]string),
	}
}

func (r *inMemoryCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
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
	r.byPost[comment.PostID] = append(r.byPost[comment.PostID], comment.ID)

	if comment.ParentID != nil {
		r.byParent[*comment.ParentID] = append(r.byParent[*comment.ParentID], comment.ID)
	}
	return nil
}

// коммент по id
func (r *inMemoryCommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, domain.ErrCommentNotFound
	}
	return comment, nil
}

// комменты по посту с пагинацией
func (r *inMemoryCommentRepository) GetByPost(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentIDs, exists := r.byPost[postID]
	if !exists {
		return []*domain.Comment{}, nil
	}

	comments := make([]*domain.Comment, 0, len(commentIDs))

	for _, id := range commentIDs {
		if comment, ok := r.comments[id]; ok {
			comments = append(comments, comment)
		}
	}

	sortCommentsByCreatedAt(comments)

	//пагинация
	start := offset
	if start > len(comments) {
		start = len(comments)
	}
	end := start + limit
	if end > len(comments) {
		end = len(comments)
	}

	return comments[start:end], nil
}

// дочерние комменты по родителю с пагинацией
func (r *inMemoryCommentRepository) GetChildren(ctx context.Context, parentID string, limit, offset int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentsIDs, exists := r.byParent[parentID]
	if !exists {
		return []*domain.Comment{}, nil
	}

	comments := make([]*domain.Comment, 0, len(commentsIDs))
	for _, id := range commentsIDs {
		if comment, ok := r.comments[id]; ok {
			comments = append(comments, comment)
		}
	}
	sortCommentsByCreatedAt(comments)

	//пагинация
	start := offset
	if start > len(comments) {
		start = len(comments)
	}
	end := start + limit
	if end > len(comments) {
		end = len(comments)
	}

	return comments[start:end], nil
}

func (r *inMemoryCommentRepository) CountByPost(ctx context.Context, postID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentIDs, exists := r.byPost[postID]
	if !exists {
		return 0, nil
	}

	return len(commentIDs), nil
}

func (r *inMemoryCommentRepository) CountChildren(ctx context.Context, parentID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commentIDs, exists := r.byParent[parentID]
	if !exists {
		return 0, nil
	}

	return len(commentIDs), nil
}

func sortCommentsByCreatedAt(comments []*domain.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
}
