package rab

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	rabService Service
}

func SetupController(r *gin.Engine, s Service, authMiddleware gin.HandlerFunc) {
	handler := &Controller{
		rabService: s,
	}

	rabGrp := r.Group("/rab")
	rabGrp.Use(authMiddleware)
	{
		rabGrp.GET("/get-rabs", handler.GetRabs)
	}
}

func (h *Controller) GetRabs(c *gin.Context) {
	q := c.Query("q")
	code := c.Query("code")

	rabs, err := h.rabService.GetListRab(c.Request.Context(), q, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rabs)
}
