package main

import (
	"mygo/config"
	"mygo/core/user"
	"mygo/modules/rab"
	"mygo/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect DB and Cache
	config.ConnectinDB()
	config.InitRedis()

	// 2. Setup Dependencies
	// Core
	userRepo := user.NewMySQLUserRepository(config.DB, config.RedisClient)
	userService := user.NewUserService(userRepo)
	
	// Modules
	rabRepo := rab.NewMySQLRabRepository(config.DB)
	rabService := rab.NewRabService(rabRepo)

	// 3. Setup Router
	r := gin.Default()
	
	// Default route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Golang Web Server is running with Modular (Package-by-Feature) Architecture!")
	})

	// Register all endpoints
	routes.SetupRoutes(r, userService, rabService)

	// 4. Run Server
	r.Run(":6000")
}
