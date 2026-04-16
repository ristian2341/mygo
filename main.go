package main

import (
	"mygo/config"
	"mygo/repository"
	"mygo/routes"
	"mygo/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect DB and Cache
	config.ConnectinDB()
	config.InitRedis()

	// 2. Setup Dependencies
	userRepo := repository.NewMySQLUserRepository(config.DB, config.RedisClient)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 3. Setup Router
	r := gin.Default()
	
	// Default route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Golang Web Server is running with Clean Architecture!")
	})

	// Register all endpoints
	routes.SetupRoutes(r, userUsecase)

	// 4. Run Server
	r.Run(":6000")
}
