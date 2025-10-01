// repositories/comment_repository.go
package repositories

import (
	"context"
	"fmt"

	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	query := `
        INSERT INTO comments (user_id, post_id, content, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, comment.UserID, comment.PostID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return &comment, nil
}
