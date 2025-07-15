// graph/resolver.go

package graph

import (
	myRedis "github.com/tmozzze/SasPosts/internal/redis"
	"github.com/tmozzze/SasPosts/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostRepo    repository.PostRepository
	CommentRepo repository.CommentRepository
	PubSub      myRedis.PubSub
}

func NewResolver(postRepo repository.PostRepository, commentRepo repository.CommentRepository, pubsub myRedis.PubSub) *Resolver {
	return &Resolver{
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
		PubSub:      pubsub,
	}
}
