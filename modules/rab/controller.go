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
		rabGrp.POST("/", handler.CreateRab)
		rabGrp.GET("/", handler.GetRabs)
	}
}

func (h *Controller) CreateRab(c *gin.Context) {
	// Dummy endpoint
	var input struct {
		ProjectName string `json:"project_name" binding:"required"`
		TotalAmount int64  `json:"total_amount" binding:"required"`
		Username    string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.rabService.CreateRab(c.Request.Context(), input.ProjectName, input.TotalAmount, input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "RAB created successfully"})
}

func (h *Controller) GetRabs(c *gin.Context) {
	rabs, err := h.rabService.GetListRab(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rabs)
}
