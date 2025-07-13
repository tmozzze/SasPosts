package graph

import (
	"github.com/tmozzze/SasPosts/internal/repository"
)

type Resolver struct {
	PostRepo    repository.PostRepository
	CommentRepo repository.CommentRepository
}
