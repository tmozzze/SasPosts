package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmozzze/SasPosts/graph"
	"github.com/tmozzze/SasPosts/graph/generated"
	"github.com/tmozzze/SasPosts/internal/config"
	"github.com/tmozzze/SasPosts/internal/repository"
	"github.com/tmozzze/SasPosts/internal/repository/inmemory"
	"github.com/tmozzze/SasPosts/internal/repository/postgres"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var postRepo repository.PostRepository
	var commentRepo repository.CommentRepository

	switch cfg.DBType {
	case "postgres":
		log.Println("Use postgres")
		dbpool, err := pgxpool.New(context.Background(), cfg.PGURL)
		if err != nil {
			log.Fatalf("failed connect to postgres %v", err)
		}
		defer dbpool.Close()

		if err := dbpool.Ping(context.Background()); err != nil {
			log.Fatalf("failed connect to postgres %v", err)
		}

		log.Println("Successfully connected")

		postRepo = postgres.NewPostgresPostRepository(dbpool)
		commentRepo = postgres.NewPostgresCommentRepository(dbpool)

	default:
		log.Println("use in-memory")
		postRepo = inmemory.NewInMemoryPostRepository()
		commentRepo = inmemory.NewInMemoryCommentRepository()
	}

	resolver := graph.NewResolver(postRepo, commentRepo)

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	server.SetErrorPresenter(graph.ErrorPresenter)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	address := ":" + cfg.Port
	log.Printf("Server on %s/ for GraphQL", cfg.Port)
	log.Fatal(http.ListenAndServe(address, nil))

}
