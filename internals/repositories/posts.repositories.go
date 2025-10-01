package repositories

import (
	"context"
	"fmt"

	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PostRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewPostRepository(db *pgxpool.Pool, rdb *redis.Client) *PostRepository {
	return &PostRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, post models.Post) (int, error) {
	query := `
        INSERT INTO posts (user_id, content, image, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id`

	var postID int

	err := r.db.QueryRow(ctx, query, post.UserID, post.Content, post.Image).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert post: %w", err)
	}

	return postID, nil
}

func (r *PostRepository) GetFeed(ctx context.Context, userID int) ([]models.PostResponse, error) {
	query := `
	SELECT p.id, p.user_id, COALESCE(u.name,'') AS name, p.content, p.image, p.created_at
	FROM posts p
	JOIN follows f ON f.following = p.user_id
	JOIN users u ON u.id = p.user_id
	WHERE f.follower = $1
	ORDER BY p.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var p models.PostResponse
		if err := rows.Scan(&p.ID, &p.UserID, &p.UserName, &p.Content, &p.Image, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
