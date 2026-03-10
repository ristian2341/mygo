package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mygo/config"
	"mygo/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomToken() string {
	// format waktu: yyyymmddhhmmss
	timeStr := time.Now().Format("20060102150405")

	// random 4 byte (8 karakter hex)
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomStr := hex.EncodeToString(randomBytes)
	rawToken := fmt.Sprintf("%s%s", timeStr, randomStr)
	hash := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(hash[:])

	// hasil akhir
	return token
}

func GenerateTokenHandler(code string, n int) string {
	date := time.Now().Format("02012006") // ddmmyyyy
	rawToken := code + date
	hash := sha256.Sum256([]byte(rawToken))
	fmt.Println(hash)
	return hex.EncodeToString(hash[:])
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	
	var dataUser models.User
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil data dari form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// query tabel user //
	err := config.DB.Where("(username = ? or email = ? )", username, username).First(&dataUser).Error
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  err.Error(),
			"status": "400",
		})
		return
	}

	// Hash password
	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(password))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Password tidak valid",
			"status": "400",
		})
		return
	}

	// generate token
	token := GenerateTokenHandler(dataUser.Code, 64)
	// simpan ke redis 12 jam
	err = config.RedisClient.Set(
		config.Ctx,
		token,
		dataUser.Username,
		12*time.Hour,
	).Err()

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal menyimpan token",
			"status": "400",
		})
		return
	}

	loginAt := time.Now().Format("2006-01-02 15:04:05")

	err = config.DB.Model(&models.User{}).
		Where("code = ?", dataUser.Code).
		Updates(map[string]interface{}{
			"login_at":     loginAt,
			"access_token": token,
		}).Error

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal Login",
			"status": "400",
		})
		return
	}

	// response
	response := map[string]interface{}{
		"username": dataUser.Username,
		"name":     dataUser.Nama,
		"token":    token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
