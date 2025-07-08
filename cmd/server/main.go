package main

import (
	"log"

	"github.com/tmozzze/SasPosts/internal/config"
	"github.com/tmozzze/SasPosts/internal/repository"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var postRepo repository.PostRepository
	var commentRepo repository.CommentRepository

	switch cfg.Database.Type {
	case "postgres":
		postRepo, commentRepo, err = repository.NewPostgresRepositories(cfg.PGURL)
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}
	default:
		postRepo = repository.NewInMemoryPostRepository()
		commentRepo = repository.NewInMemoryCommentRepository()
	}
}
