package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/handlers"
	"github.com/habibmrizki/finalphase3/internals/middleware"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitLikesRouter(router *gin.Engine, db *pgxpool.Pool) {
	likeRepo := repositories.NewLikeRepository(db)
	likeHandler := handlers.NewLikeHandler(likeRepo)

	likes := router.Group("/likes", middleware.VerifyToken())
	likes.POST("/:id", likeHandler.LikePost)
}
