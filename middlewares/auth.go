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

		// cek token di Redis
		_, err := config.RedisClient.Get(
			config.Ctx,
			token,
		).Result()

		if err != nil {

			// jika tidak ada di redis, cek ke database
			var user models.User
			db := config.DB

			errDB := db.Where("access_token = ?", token).First(&user).Error
			if errDB != nil {

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token invalid or expired",
				})
				c.Abort()
				return
			}

			// jika ditemukan di DB, simpan lagi ke Redis (optional)
			config.RedisClient.Set(
				config.Ctx,
				token,
				user.Username,
				12*time.Hour,
			)

		}

		c.Next()
	}
}

func GenerateResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Method not allowed",
			"status": "400",
		})
		return
	}

	// WAJIB parse form dulu
	if err := r.ParseForm(); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Gagal membaca form",
			"status": "400",
		})
		return
	}

	// Ambil data dari form
	email := r.FormValue("email")
	if email == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Email harus di isi",
			"status": "400",
		})
		return
	}

	// cek email ada di db dan generate code 6 digit //
	var dataUser models.User
	err := config.DB.Where("email = ? ", email).First(&dataUser).Error
	if err != nil {
		if err.Error() == "record not found" {
			json.NewEncoder(w).Encode(map[string]string{
				"error":  "Email : " + email + " tidak ditemukan",
				"status": "400",
			})
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
