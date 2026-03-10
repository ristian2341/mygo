package middlewares

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mygo/config"
	"mygo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token not found",
			})
			c.Abort()
			return
		}

		// nanti bisa kamu cek ke Redis
		_, err := config.RedisClient.Get(
			config.Ctx,
			token,
		).Result()

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalid or expired",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GenerateResetPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// WAJIB parse form dulu
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	// Ambil data dari form
	email := r.FormValue("email")

	if email == "" {
		http.Error(w, "Email harus diisi", http.StatusBadRequest)
		return
	}

	// cek email ada di db dan generate code 6 digit //
	var dataUser models.User
	err := config.DB.Where("email = ? ", email).First(&dataUser).Error
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Email : "+email+" tidak ditemukan", http.StatusBadRequest)
			return
		}
	}

	verify_code := Generate6DigitCode()

	// Update db User update verify_code //
	result := config.DB.Model(&models.User{}).Where("email = ? ", email).Update("verify_code", verify_code)

	// Simpan ke database
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// kirim email di middlewares
	body := "Hi, " + dataUser.Nama + "\n\n" +
		"Berikut kami kirimkan kode verifikasi untuk reset password anda.\n\n" +
		"Kode Verifikasi: " + verify_code + "\n\n" +
		"Terima kasih."
	send := SendEmail(verify_code, email, body)
	if send == false {
		http.Error(w, "Kirim Email verifikasi gagal", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kode Verifikasi Reset Password sudah terkirim Email, cek email anda atau cek di folder spam",
	})

}

func Generate6DigitCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}
