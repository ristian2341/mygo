package routes

import (
	"mygo/config"
	"mygo/controllers"
	"mygo/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/token_random", func(c *gin.Context) {
		token := controllers.GenerateRandomToken()
		// save to redis 60 menit //
		err := config.RedisClient.Set(
			config.Ctx,
			token,
			"valid",
			60*time.Minute,
		).Err()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(200, gin.H{
			"token": token, // ✅ panggil function
		})
	})

	// route wajib token
	auth := r.Group("/")
	auth.Use(middlewares.TokenRequired())

	auth.POST("/user/register-user", func(c *gin.Context) {
		controllers.RegisterUser(c.Writer, c.Request)
	})
	auth.POST("/auth/login", func(c *gin.Context) {
		controllers.GetLogin(c.Writer, c.Request)
	})
	auth.POST("/auth/reset-password", func(c *gin.Context) {
		middlewares.GenerateResetPassword(c.Writer, c.Request)
	})
	auth.POST("/user/password-reset", func(c *gin.Context) {
		controllers.PasswordReset(c.Writer, c.Request)
	})
	auth.POST("/user/ubah-reset", func(c *gin.Context) {
		controllers.ChangePasword(c.Writer, c.Request)
	})
	auth.POST("/user/update-user", func(c *gin.Context) {
		controllers.UpdateDataUser(c.Writer, c.Request)
	})
}
