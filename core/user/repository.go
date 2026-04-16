package user

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewMySQLUserRepository(db *gorm.DB, redisClient *redis.Client) Repository {
	return &mysqlUserRepository{
		db:    db,
		redis: redisClient,
	}
}

func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) GetByUsernameOrEmail(ctx context.Context, identifier string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("username = ? OR email = ?", identifier, identifier).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) GetByCodeAndUsername(ctx context.Context, code string, username string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("code = ? AND username = ?", code, username).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) GetByVerifyCode(ctx context.Context, verifyCode string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("verify_code = ?", verifyCode).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) GetByAccessToken(ctx context.Context, token string) (User, error) {
	var userData User
	err := m.db.WithContext(ctx).Where("access_token = ?", token).First(&userData).Error
	return userData, err
}

func (m *mysqlUserRepository) Create(ctx context.Context, userData *User) error {
	return m.db.WithContext(ctx).Create(userData).Error
}

func (m *mysqlUserRepository) UpdateByCode(ctx context.Context, code string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&User{}).Where("code = ?", code).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByEmail(ctx context.Context, email string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByVerifyCode(ctx context.Context, verifyCode string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&User{}).Where("verify_code = ?", verifyCode).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByUsername(ctx context.Context, username string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Updates(updates).Error
}

func (m *mysqlUserRepository) GetLastCodeByDate(ctx context.Context, todayStr string) (string, error) {
	var userData User
	err := m.db.WithContext(ctx).
		Where("code LIKE ?", todayStr+"%").
		Order("code DESC").
		First(&userData).Error

	if err != nil {
		return "", err
	}
	return userData.Code, nil
}

func (m *mysqlUserRepository) FetchAll(ctx context.Context, q string, code string) ([]User, error) {
	var users []User
	query := m.db.WithContext(ctx).Model(&User{})

	if q != "" {
		searchTerm := "%" + q + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nama LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if code != "" {
		query = query.Where("code = ?", code)
	}

	err := query.Find(&users).Error
	return users, err
}

func (m *mysqlUserRepository) SetRedisToken(ctx context.Context, token string, value string, ttl uint) error {
	duration := time.Duration(ttl) * time.Hour
	return m.redis.Set(ctx, token, value, duration).Err()
}

func (m *mysqlUserRepository) GetRedisToken(ctx context.Context, token string) (string, error) {
	return m.redis.Get(ctx, token).Result()
}
