package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	mr "math/rand"
	"mygo/domain"
	"mygo/utils"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u *userUsecase) Register(ctx context.Context, username, password, email, nama string) error {
	// cek duplikasi username
	_, err := u.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return fmt.Errorf("Username %s sudah ada, masukan username yang lain", username)
	}

	// cek duplikasi email
	_, err = u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("Email %s sudah ada, masukan email yang lain", email)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal memproses password")
	}

	// generate code
	code, err := u.generateUserCode(ctx)
	if err != nil {
		return errors.New("Gagal generate code")
	}

	user := &domain.User{
		Code:     code,
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nama:     nama,
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) generateUserCode(ctx context.Context) (string, error) {
	today := time.Now().Format("20060102")
	lastCode, err := u.userRepo.GetLastCodeByDate(ctx, today)

	sequence := 1
	if err == nil && len(lastCode) >= 12 {
		lastNumberStr := lastCode[8:] // ambil 4 digit terakhir
		if lastNumber, errConv := strconv.Atoi(lastNumberStr); errConv == nil {
			sequence = lastNumber + 1
		}
	}

	newCode := fmt.Sprintf("%s%04d", today, sequence)
	return newCode, nil
}

func (u *userUsecase) CheckToken(ctx context.Context, token string) error {
	_, err := u.userRepo.GetRedisToken(ctx, token)
	if err != nil {
		// jika tidak ada di redis, cek db
		user, errDb := u.userRepo.GetByAccessToken(ctx, token)
		if errDb != nil {
			return errors.New("Token invalid or expired")
		}
		// resave optional
		u.userRepo.SetRedisToken(ctx, token, user.Username, 12)
	}
	return nil
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (map[string]interface{}, error) {
	user, err := u.userRepo.GetByUsernameOrEmail(ctx, username)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// Hash password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Password tidak valid")
	}

	// generate token
	date := time.Now().Format("02012006")
	rawToken := user.Code + date
	hash := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(hash[:])

	// save to redis 12 hours
	err = u.userRepo.SetRedisToken(ctx, token, user.Username, 12)
	if err != nil {
		return nil, errors.New("Gagal menyimpan token")
	}

	loginAt := time.Now().Format("2006-01-02 15:04:05")

	err = u.userRepo.UpdateByCode(ctx, user.Code, map[string]interface{}{
		"login_at":     loginAt,
		"access_token": token,
	})
	if err != nil {
		return nil, errors.New("Gagal update data login")
	}

	response := map[string]interface{}{
		"username": user.Username,
		"name":     user.Nama,
		"token":    token,
	}

	return response, nil
}

func (u *userUsecase) PasswordResetRequest(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("Email : %s tidak ditemukan", email)
	}

	mr.Seed(time.Now().UnixNano())
	verifyCode := fmt.Sprintf("%06d", mr.Intn(900000)+100000)

	err = u.userRepo.UpdateByEmail(ctx, email, map[string]interface{}{
		"verify_code": verifyCode,
	})
	if err != nil {
		return errors.New("Gagal mengupdate verify code")
	}

	body := "Hi, " + user.Nama + "\n\n" +
		"Berikut kami kirimkan kode verifikasi untuk reset password anda.\n\n" +
		"Kode Verifikasi: " + verifyCode + "\n\n" +
		"Terima kasih."
	success := utils.SendEmail(verifyCode, email, body)
	if !success {
		return errors.New("Kirim Email verifikasi gagal")
	}
	return nil
}

func (u *userUsecase) PasswordResetSubmit(ctx context.Context, verifyCode, password, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("Password dan Retype Password tidak sama")
	}

	_, err := u.userRepo.GetByVerifyCode(ctx, verifyCode)
	if err != nil {
		return errors.New("Kode verifikasi tidak ditemukan")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal memproses kode password")
	}

	return u.userRepo.UpdateByVerifyCode(ctx, verifyCode, map[string]interface{}{
		"password":    string(hashPassword),
		"verify_code": "",
	})
}

func (u *userUsecase) ChangePassword(ctx context.Context, username, oldPassword, newPassword, passwordConfirm string) error {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return errors.New("Username tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("Password lama tidak valid")
	}

	if newPassword != passwordConfirm {
		return errors.New("Password baru dan konfirmasi tidak sama")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal memproses password")
	}

	return u.userRepo.UpdateByUsername(ctx, username, map[string]interface{}{
		"password": string(hashPassword),
	})
}

func (u *userUsecase) UpdateProfile(ctx context.Context, code, username, email, nama, foto, phone, supervisor string) error {
	_, err := u.userRepo.GetByCodeAndUsername(ctx, code, username)
	if err != nil {
		return errors.New("Data user tidak ditemukan")
	}

	// cek duplikasi email lain
	userEmail, errEmail := u.userRepo.GetByEmail(ctx, email)
	if errEmail == nil && userEmail.Code != code {
		return errors.New("Email sudah dipakai oleh user lain")
	}

	return u.userRepo.UpdateByUsername(ctx, username, map[string]interface{}{
		"email":      email,
		"nama":       nama,
		"foto":       foto,
		"phone":      phone,
		"supervisor": supervisor,
	})
}

func (u *userUsecase) GetUsers(ctx context.Context, q, code string) ([]domain.User, error) {
	return u.userRepo.FetchAll(ctx, q, code)
}

func (u *userUsecase) GenerateRandomToken(ctx context.Context) (string, error) {
	timeStr := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomStr := hex.EncodeToString(randomBytes)
	rawToken := fmt.Sprintf("%s%s", timeStr, randomStr)
	hash := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(hash[:])
	return token, nil
}
