package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) LikePost(ctx context.Context, userID, postID int) error {
	query := `INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}
	return nil
}
