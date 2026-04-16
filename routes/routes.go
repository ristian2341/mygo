package routes

import (
	"mygo/domain"
	"mygo/handler"
	"mygo/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, uc domain.UserUsecase) {
	// middlewares
	authMw := middlewares.TokenRequired(uc)

	// endpoints registration
	handler.NewAuthHandler(r, uc)
	handler.NewUserHandler(r, uc, authMw)
}
