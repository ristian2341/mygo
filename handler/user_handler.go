package handler

import (
	"mygo/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, us domain.UserUsecase, authMiddleware gin.HandlerFunc) {
	handler := &UserHandler{
		userUsecase: us,
	}

	userGrp := r.Group("/user")
	userGrp.Use(authMiddleware)
	{
		userGrp.POST("/register-user", handler.RegisterUser)
		userGrp.POST("/password-reset", handler.PasswordResetSubmit)
		userGrp.POST("/ubah-reset", handler.ChangePassword)
		userGrp.POST("/update-user", handler.UpdateUser)
		userGrp.GET("/get-users", handler.GetUsers)
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	nama := c.PostForm("nama")

	if username == "" || password == "" || email == "" || nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua Field harus diisi", "status": "400"})
		return
	}

	err := h.userUsecase.Register(c.Request.Context(), username, password, email, nama)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil didaftarkan",
	})
}

func (h *UserHandler) PasswordResetSubmit(c *gin.Context) {
	verifyCode := c.PostForm("verify_code")
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password_confirm")

	if verifyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode Verikasi tidak boleh kosong", "status": "400"})
		return
	}
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak boleh kosong", "status": "400"})
		return
	}
	if passwordConfirm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Retype Password tidak boleh kosong", "status": "400"})
		return
	}

	err := h.userUsecase.PasswordResetSubmit(c.Request.Context(), verifyCode, password, passwordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // using default format from old code
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset Password Berhasil."})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	username := c.PostForm("username")
	passwordLama := c.PostForm("password_lama")
	passwordBaru := c.PostForm("password_baru")
	passwordConfirm := c.PostForm("password_confirm")

	if passwordLama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak boleh kosong", "status": "400"})
		return
	}
	if passwordConfirm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Retype Password tidak boleh kosong", "status": "400"})
		return
	}

	err := h.userUsecase.ChangePassword(c.Request.Context(), username, passwordLama, passwordBaru, passwordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ubah Password Berhasil."})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	code := c.PostForm("code")
	username := c.PostForm("username")
	email := c.PostForm("email")
	nama := c.PostForm("nama")
	foto := c.PostForm("foto")
	phone := c.PostForm("phone")
	supervisor := c.PostForm("supervisor")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak boleh kosong", "status": "400"})
		return
	}

	err := h.userUsecase.UpdateProfile(c.Request.Context(), code, username, email, nama, foto, phone, supervisor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update User Berhasil."})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	q := c.Query("q")
	code := c.Query("code")

	users, err := h.userUsecase.GetUsers(c.Request.Context(), q, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
