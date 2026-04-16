package handler

import (
	"mygo/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUsecase domain.UserUsecase
}

func NewAuthHandler(r *gin.Engine, us domain.UserUsecase) {
	handler := &AuthHandler{
		userUsecase: us,
	}

	r.GET("/token_random", handler.GenerateToken)
	
	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/reset-password", handler.ResetPasswordRequest)
	}
}

func (h *AuthHandler) GenerateToken(c *gin.Context) {
	token, err := h.userUsecase.GenerateRandomToken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// Because original code uses r.FormValue, we bind form data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "status": "400"})
		return
	}

	res, err := h.userUsecase.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) ResetPasswordRequest(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email harus di isi", "status": "400"})
		return
	}

	err := h.userUsecase.PasswordResetRequest(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kode Verifikasi Reset Password sudah terkirim Email, cek email anda atau cek di folder spam",
	})
}
