package main

import (
	"log"
	"os"

	"github.com/habibmrizki/finalphase3/internals/configs"
	"github.com/habibmrizki/finalphase3/internals/routers"
	"github.com/joho/godotenv"
)

// @title 		Final Phase 3
// @version		1.0
// @description	RESTful API created using gin for Final Phase 3
// @host		localhost:3000
// @basepath	/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env\nCause", err.Error())
		return
	}
	log.Println(os.Getenv("DBUSER"))

	db, err := configs.InitDb()
	if err != nil {
		log.Println("Failed to connect to database\nCause", err.Error())
		return
	}
	log.Println("DB CONNECTED")

	defer db.Close()

	rdb, err := configs.InitRedis()

	if err != nil {
		log.Println("Failed to connect to Redis\nCause", err.Error())
		return
	}
	log.Println("Redis Connected")

	router := routers.InitRouter(db, rdb)
	router.Run(":3000")
}
