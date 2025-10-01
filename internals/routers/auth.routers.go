package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/handlers"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/habibmrizki/finalphase3/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitUsersRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	log.Println("Initializing Users/Auth Routes...")

	hashConfig := pkg.NewHashConfig()

	authRepo := repositories.NewAuthRepository(db, rdb)

	authHandler := handlers.NewAuthHandler(authRepo, hashConfig)

	usersGroup := router.Group("/auth")

	usersGroup.POST("/register", authHandler.Register)
	usersGroup.POST("/login", authHandler.Login)

}
