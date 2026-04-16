package repository

import (
	"context"
	"mygo/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewMySQLUserRepository(db *gorm.DB, redisClient *redis.Client) domain.UserRepository {
	return &mysqlUserRepository{
		db:    db,
		redis: redisClient,
	}
}

func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) GetByUsernameOrEmail(ctx context.Context, identifier string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("username = ? OR email = ?", identifier, identifier).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) GetByCodeAndUsername(ctx context.Context, code string, username string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("code = ? AND username = ?", code, username).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) GetByVerifyCode(ctx context.Context, verifyCode string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("verify_code = ?", verifyCode).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) GetByAccessToken(ctx context.Context, token string) (domain.User, error) {
	var user domain.User
	err := m.db.WithContext(ctx).Where("access_token = ?", token).First(&user).Error
	return user, err
}

func (m *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *mysqlUserRepository) UpdateByCode(ctx context.Context, code string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&domain.User{}).Where("code = ?", code).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByEmail(ctx context.Context, email string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByVerifyCode(ctx context.Context, verifyCode string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&domain.User{}).Where("verify_code = ?", verifyCode).Updates(updates).Error
}

func (m *mysqlUserRepository) UpdateByUsername(ctx context.Context, username string, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&domain.User{}).Where("username = ?", username).Updates(updates).Error
}

func (m *mysqlUserRepository) GetLastCodeByDate(ctx context.Context, todayStr string) (string, error) {
	var user domain.User
	err := m.db.WithContext(ctx).
		Where("code LIKE ?", todayStr+"%").
		Order("code DESC").
		First(&user).Error

	if err != nil {
		return "", err
	}
	return user.Code, nil
}

func (m *mysqlUserRepository) FetchAll(ctx context.Context, q string, code string) ([]domain.User, error) {
	var users []domain.User
	query := m.db.WithContext(ctx).Model(&domain.User{})

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
