package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/handlers"
	"github.com/habibmrizki/finalphase3/internals/middleware"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	// Import middleware JWT Anda
	// "github.com/habibmrizki/finalphase3/internals/middleware"
)

func InitPostsRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	postRepo := repositories.NewPostRepository(db, rdb)
	postHandler := handlers.NewPostHandler(postRepo, rdb)

	postsGroup := router.Group("/posts", middleware.VerifyToken())

	postsGroup.POST("/", postHandler.CreatePost)
	postsGroup.GET("/feed", postHandler.GetFeed)
}
