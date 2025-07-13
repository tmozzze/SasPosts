package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmozzze/SasPosts/internal/domain"
	"github.com/tmozzze/SasPosts/utils"
)

type PostgresPostRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostRepository(db *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) Create(ctx context.Context, post *domain.Post) error {
	if post.ID == "" {
		post.ID = utils.GenerateID()
	}
	post.CreatedAt = time.Now()

	query := `INSERT INTO posts (id, title, content, author, allow_comments, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query,
		post.ID,
		post.Title,
		post.Content,
		post.Author,
		post.AllowComments,
		post.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *PostgresPostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	query := `SELECT id, title, content, author, allow_comments, created_at 
	          FROM posts WHERE id = $1`
	var post domain.Post
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Author,
		&post.AllowComments,
		&post.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("id search post failed: %w", err)
	}

	return &post, nil

}

func (r *PostgresPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := `SELECT id, title, content, author, allow_comments, created_at FROM posts`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed get all posts %w", err)
	}
	defer rows.Close()

	var posts []*domain.Post

	for rows.Next() {
		var post domain.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&post.AllowComments,
			&post.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed scan posts %w", err)
		}

		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error %w", err)
	}

	return posts, nil
}

func (r *PostgresPostRepository) ToggleComments(ctx context.Context, postID string, allow bool) error {
	query := `UPDATE posts SET allow_comments = $1 WHERE id = $2`

	commandTag, err := r.db.Exec(ctx, query, allow, postID)
	if err != nil {
		return fmt.Errorf("failed toggle comments %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}

func (r *PostgresPostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, author = $3, allow_comments = $4
			  WHERE id = $5`

	commantTag, err := r.db.Exec(ctx, query,
		post.Title,
		post.Content,
		post.Author,
		post.AllowComments,
		post.ID,
	)

	if err != nil {
		return fmt.Errorf("failed update post %w", err)
	}

	if commantTag.RowsAffected() == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}

func (r *PostgresPostRepository) Delete(ctx context.Context, postID string) error {
	query := `DELETE FROM posts WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, postID)

	if err != nil {
		return fmt.Errorf("failed delete post %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}

func (r *PostgresPostRepository) CheckAllowedComments(ctx context.Context, postID string) (bool, error) {
	query := `SELECT allow_comments FROM posts WHERE id = $1`

	var allowComments bool
	err := r.db.QueryRow(ctx, query, postID).Scan(&allowComments)

	if err != nil {
		return false, fmt.Errorf("failed check allow comments %w", err)
	}

	return allowComments, nil
}
