package controllers

import (
	"encoding/json"
	"fmt"
	"mygo/config"
	"mygo/models"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
			"status": "400",
		})
		return
	}

	// WAJIB parse form dulu
	if err := r.ParseForm(); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal membaca form",
			"status": "400",
		})
		return
	}

	// Ambil data dari form
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	nama := r.FormValue("nama")

	// Validasi
	if username == "" || password == "" || email == "" || nama == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Semua Field harus diisi",
			"status": "400",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal memproses password",
			"status": "400",
		})
		return
	}

	var dataUser models.User

	// cek user name //
	err = config.DB.Where("username = ?", username).First(&dataUser).Error
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Username "+username+"sudah ada, masukan email yang lain",
			"status": "400",
		})
		return
	}

	// cek email //
	err = config.DB.Where("email = ?", email).First(&dataUser).Error
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email "+email+" sudah ada, masukan email yang lain",
			"status": "400",
		})
		return
	}

	code, err := GenerateUserCode()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal generate code",
			"status": "400",
		})
		return
	}

	// Masukkan ke struct
	user := models.User{
		Code:     code,
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nama:     nama,
	}

	// Simpan ke database
	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User berhasil didaftarkan",
	})
}

func GenerateUserCode() (string, error) {

	var lastUser models.User
	today := time.Now().Format("20060102")

	// cari kode terakhir hari ini
	err := config.DB.
		Where("code LIKE ?", today+"%").
		Order("code DESC").
		First(&lastUser).Error

	sequence := 1

	if err == nil {
		lastCode := lastUser.Code
		lastNumberStr := lastCode[8:] // ambil 4 digit terakhir
		lastNumber, _ := strconv.Atoi(lastNumberStr)
		sequence = lastNumber + 1
	}

	newCode := fmt.Sprintf("%s%04d", today, sequence)

	return newCode, nil
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
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
	verify_code := r.FormValue("verify_code")
	password := r.FormValue("password")
	password_confirm := r.FormValue("password_confirm")

	if verify_code == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Kode Verikasi tidak boleh kosong",
			"status": "400",
		})
		return
	}

	if password == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Password tidak boleh kosong",
			"status": "400",
		})
		return
	}

	if password_confirm == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Retype Password tidak boleh kosong",
			"status": "400",
		})
		return
	}

	if password_confirm != password {
		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Password dan Retype Password tidak sama",
			"status": "400",
		})
		return
	}

	// data user //
	var dataUser models.User
	err := config.DB.Where("verify_code = ? ", verify_code).First(&dataUser).Error
	if err != nil {
		if err.Error() == "record not found" {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Kode verifikasi tidak ditemukan, Mohon cek kembali email anda",
			})
			return
		}
	}

	// Hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal memproses password",
		})
		return
	}

	// Update db User update verify_code //
	result := config.DB.Model(&models.User{}).Where("verify_code = ? ", verify_code).Update("password", hashPassword)

	if result.Error != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Gagal memproses password",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reset Password Berhasil.",
	})
}
