package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	userService Service
}

func SetupController(r *gin.Engine, us Service, authMiddleware gin.HandlerFunc) {
	handler := &Controller{
		userService: us,
	}

	// Public Auth Endpoints
	r.GET("/token_random", handler.GenerateToken)
	
	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/reset-password", handler.ResetPasswordRequest)
	}

	// Protected User Endpoints
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

func (h *Controller) GenerateToken(c *gin.Context) {
	token, err := h.userService.GenerateRandomToken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Controller) Login(c *gin.Context) {
	var input struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "status": "400"})
		return
	}

	res, err := h.userService.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Controller) ResetPasswordRequest(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email harus di isi", "status": "400"})
		return
	}

	err := h.userService.PasswordResetRequest(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kode Verifikasi Reset Password sudah terkirim Email, cek email anda atau cek di folder spam",
	})
}

func (h *Controller) RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	nama := c.PostForm("nama")

	if username == "" || password == "" || email == "" || nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua Field harus diisi", "status": "400"})
		return
	}

	err := h.userService.Register(c.Request.Context(), username, password, email, nama)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil didaftarkan",
	})
}

func (h *Controller) PasswordResetSubmit(c *gin.Context) {
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

	err := h.userService.PasswordResetSubmit(c.Request.Context(), verifyCode, password, passwordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset Password Berhasil."})
}

func (h *Controller) ChangePassword(c *gin.Context) {
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

	err := h.userService.ChangePassword(c.Request.Context(), username, passwordLama, passwordBaru, passwordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "400"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ubah Password Berhasil."})
}

func (h *Controller) UpdateUser(c *gin.Context) {
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

	err := h.userService.UpdateProfile(c.Request.Context(), code, username, email, nama, foto, phone, supervisor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update User Berhasil."})
}

func (h *Controller) GetUsers(c *gin.Context) {
	q := c.Query("q")
	code := c.Query("code")

	users, err := h.userService.GetUsers(c.Request.Context(), q, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
