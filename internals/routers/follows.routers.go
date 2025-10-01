package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/handlers"
	"github.com/habibmrizki/finalphase3/internals/middleware"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitFollowRouter(router *gin.Engine, db *pgxpool.Pool) {
	followRepo := repositories.NewFollowRepository(db)
	followHandler := handlers.NewFollowHandler(followRepo)

	followRoutes := router.Group("/api/follows", middleware.VerifyToken())

	followRoutes.POST("/:id", followHandler.Follow)
}
