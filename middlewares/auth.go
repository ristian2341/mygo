package middlewares

import (
	"mygo/core/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenRequired(uc user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token not found",
			})
			c.Abort()
			return
		}

		err := uc.CheckToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
