package user

import "context"

type User struct {
	Code        string  `json:"code" binding:"required,len=14"`
	Username    string  `json:"username" binding:"required,min=3,max=100"`
	Password    string  `json:"-" binding:"required,min=6,max=150"`
	Email       string  `json:"email" binding:"required,email"`
	Nama        string  `json:"nama" binding:"required,max=150"`
	Phone       string  `json:"phone"`
	Foto        string  `json:"foto"`
	Supervisor  *uint   `json:"supervisor"`
	Status      *uint   `json:"status"`
	VerifyCode  *string `json:"verify_code,omitempty"`
	LoginAt     *string `json:"login_at,omitempty"`
	AccessToken *string `json:"access_token,omitempty"`
}

// TableName meng-override nama tabel default gorm
func (User) TableName() string {
	return "user" // isi dengan nama tabel yang kamu mau
}

type Repository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsernameOrEmail(ctx context.Context, identifier string) (User, error)
	GetByCodeAndUsername(ctx context.Context, code string, username string) (User, error)
	GetByVerifyCode(ctx context.Context, verifyCode string) (User, error)
	GetByAccessToken(ctx context.Context, token string) (User, error)
	Create(ctx context.Context, user *User) error
	UpdateByCode(ctx context.Context, code string, updates map[string]interface{}) error
	UpdateByEmail(ctx context.Context, email string, updates map[string]interface{}) error
	UpdateByVerifyCode(ctx context.Context, verifyCode string, updates map[string]interface{}) error
	UpdateByUsername(ctx context.Context, username string, updates map[string]interface{}) error
	GetLastCodeByDate(ctx context.Context, todayStr string) (string, error)
	FetchAll(ctx context.Context, q string, code string) ([]User, error)

	// Redis related ops 
	SetRedisToken(ctx context.Context, token string, value string, ttl uint) error
	GetRedisToken(ctx context.Context, token string) (string, error)
}

type Service interface {
	Register(ctx context.Context, username, password, email, nama string) error
	Login(ctx context.Context, username, password string) (map[string]interface{}, error)
	PasswordResetRequest(ctx context.Context, email string) error
	PasswordResetSubmit(ctx context.Context, verifyCode, password, passwordConfirm string) error
	ChangePassword(ctx context.Context, username, oldPassword, newPassword, passwordConfirm string) error
	UpdateProfile(ctx context.Context, code, username, email, nama, foto, phone, supervisor string) error
	GetUsers(ctx context.Context, q, code string) ([]User, error)
	CheckToken(ctx context.Context, token string) (string, error)
	GenerateRandomToken(ctx context.Context) (string, error)
}
