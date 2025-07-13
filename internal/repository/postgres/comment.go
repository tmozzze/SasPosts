package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmozzze/SasPosts/internal/domain"
)

type PostgresCommentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCommentRepository(db *pgxpool.Pool) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	if utf8.RuneCountInString(comment.Content) > domain.MaxCommentLength {
		return domain.ErrCommentTooLong
	}

	if comment.ParentID == nil {
		comment.Depth = 0
		comment.Path = comment.ID
	} else {
		var parentPath string
		var parentDepth int

		query := `SELECT path, depth FROM comments WHERE id = $1`
		err := r.db.QueryRow(ctx, query, *comment.ParentID).Scan(&parentPath, &parentDepth)
		if err != nil {
			if err == pgx.ErrNoRows {
				return domain.ErrParentCommentNotFound
			}
			return fmt.Errorf("failed get parent comment %w", err)
		}

		comment.Depth = parentDepth + 1
		comment.Path = parentPath + "." + comment.ID
	}

	insertQuery := `INSERT INTO comments (id, post_id, parent_id, author, content, path, depth, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(ctx, insertQuery,
		comment.ID,
		comment.PostID,
		comment.ParentID,
		comment.Author,
		comment.Content,
		comment.Path,
		comment.Depth,
		comment.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed create comment %w", err)
	}

	return nil
}

func (r *PostgresCommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	query := `SELECT id, post_id, parent_id, author, content, path, depth, created_at
			  FROM comments WHERE id = $1`

	var comment domain.Comment
	var scannedParentID sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&scannedParentID,
		&comment.Author,
		&comment.Content,
		&comment.Path,
		&comment.Depth,
		&comment.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed get comment by id %w", err)
	}

	if scannedParentID.Valid {
		comment.ParentID = &scannedParentID.String
	}

	return &comment, nil
}

func (r *PostgresCommentRepository) GetByPost(ctx context.Context, postID string, limit int, offset int) ([]*domain.Comment, error) {
	query := `SELECT id, post_id, parent_id, author, content, path, depth, created_at
			  FROM comments WHERE post_id = $1 AND parent_id IS NULL
			  ORDER BY created_at ASC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, postID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed get comments by post %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment

	for rows.Next() {
		var comment domain.Comment
		var scannedParentID sql.NullString

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&scannedParentID,
			&comment.Author,
			&comment.Content,
			&comment.Path,
			&comment.Depth,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed scan comment %w", err)
		}

		if scannedParentID.Valid {
			comment.ParentID = &scannedParentID.String
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *PostgresCommentRepository) GetChildren(ctx context.Context, parentID string, limit int, offset int) ([]*domain.Comment, error) {
	query := `SELECT id, post_id, parent_id, author, content, path, depth, created_at
			  FROM comments WHERE parent_id = $1
			  ORDER BY created_at ASC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, parentID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed get comments by post %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment

	for rows.Next() {
		var comment domain.Comment
		var scannedParentID sql.NullString

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&scannedParentID,
			&comment.Author,
			&comment.Content,
			&comment.Path,
			&comment.Depth,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed scan comment %w", err)
		}

		if scannedParentID.Valid {
			comment.ParentID = &scannedParentID.String
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *PostgresCommentRepository) CountByPost(ctx context.Context, postID string) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int

	err := r.db.QueryRow(ctx, query, postID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed count comments by post %w", err)
	}

	return count, nil
}

func (r *PostgresCommentRepository) CountChildren(ctx context.Context, parentID string) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE parent_id = $1`

	var count int

	err := r.db.QueryRow(ctx, query, parentID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed count comments by children %w", err)
	}
	return count, nil
}
