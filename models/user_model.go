package models

type User struct {
	Code     string `json:"code" binding:"required,len=14"`
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=150"`
	Email    string `json:"email" binding:"required,email"`
	Nama     string `json:"nama" binding:"required,max=150"`
}

// Override nama tabel
func (User) TableName() string {
	return "user" // isi dengan nama tabel yang kamu mau
}
