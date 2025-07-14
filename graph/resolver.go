// graph/resolver.go

package graph

import (
	"github.com/redis/go-redis/v9"
	"github.com/tmozzze/SasPosts/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostRepo    repository.PostRepository
	CommentRepo repository.CommentRepository
	Redis       *redis.Client
}

func NewResolver(postRepo repository.PostRepository, commentRepo repository.CommentRepository, redisClient *redis.Client) *Resolver {
	return &Resolver{
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
		Redis:       redisClient,
	}
}
