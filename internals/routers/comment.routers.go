package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/handlers"
	"github.com/habibmrizki/finalphase3/internals/middleware"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCommentsRouter(router *gin.Engine, db *pgxpool.Pool) {
	commentRepo := repositories.NewCommentRepository(db)
	commentHandler := handlers.NewCommentHandler(commentRepo)

	comments := router.Group("/comments", middleware.VerifyToken())

	comments.POST("/:id", commentHandler.CreateComment)
}
