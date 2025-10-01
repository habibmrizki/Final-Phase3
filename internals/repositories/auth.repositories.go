package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Ganti AuhtRepository menjadi AuthRepository
type AuthRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, name, avatar, bio FROM users WHERE email = $1"
	user := models.User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.Bio,
	)

	if err != nil {
		log.Println("FindByEmail error:", err)
		return nil, err
	}
	return &user, nil
}

// CreateUser mendaftarkan user baru
func (r *AuthRepository) CreateUser(ctx context.Context, req models.RegisterRequest, hashedPassword string) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, password, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, email, name`

	newUser := models.User{}

	err := r.db.QueryRow(ctx, query, req.Name, req.Email, hashedPassword).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	// Langsung kembalikan struct yang sudah di-scan
	return &newUser, nil
}
