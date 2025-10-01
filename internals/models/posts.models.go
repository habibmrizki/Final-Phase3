package models

import (
	"mime/multipart"
	"time"
)

type CreatePostRequest struct {
	Content string                `form:"content"`
	Image   *multipart.FileHeader `form:"image"`
}

type Post struct {
	ID        int
	UserID    int
	Content   string
	Image     string
	CreatedAt time.Time
}

type PostResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}
