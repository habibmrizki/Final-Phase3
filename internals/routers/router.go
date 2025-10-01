package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/docs"
	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	log.Println("Router created.")

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Static("/uploads", "./public/uploads")

	InitPostsRouter(router, db, rdb)
	InitLikesRouter(router, db)
	InitCommentsRouter(router, db)
	InitUsersRouter(router, db, rdb)
	InitFollowRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			Message: "Rute Salah",
			Status:  "Rute Tidak Ditemukan",
		})
	})

	return router
}
