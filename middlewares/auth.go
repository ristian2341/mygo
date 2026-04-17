package middlewares

import (
	"mygo/core/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenRequired(uc user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// Buang teks "Bearer " jika disertakan oleh Postman/Klien
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Anda belum Login. Silakan menyertakan token akses terlebih dahulu.",
			})
			c.Abort()
			return
		}

		userCode, err := uc.CheckToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_code", userCode)
		c.Next()
	}
}
