package routes

import (
	"mygo/core/user"
	"mygo/middlewares"
	"mygo/modules/rab"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, us user.Service, rabSvc rab.Service) {
	// middlewares
	authMw := middlewares.TokenRequired(us)

	// endpoints registration
	user.SetupController(r, us, authMw)
	rab.SetupController(r, rabSvc, authMw)
}
