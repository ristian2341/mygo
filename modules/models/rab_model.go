package models

import (
	"time"

	"gorm.io/gorm"
)

type Rab struct {
	ID          uint           `gorm:"primaryKey"`
	Code        string         `json:"code"`
	CodeProject string         `json:"code_project"`
	TglMulai    time.Time      `json:"tgl_mulai"`
	TglAkhir    time.Time      `json:"tgl_akhir"`

	CreatedBy uint `json:"created_by"`
	UpdatedBy uint `json:"updated_by"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Rab) BeforeCreate(tx *gorm.DB) (err error) {
	userID := GetUserID(tx)

	r.CreatedBy = userID
	r.UpdatedBy = userID

	return
}

func (r *Rab) BeforeUpdate(tx *gorm.DB) (err error) {
	userID := GetUserID(tx)

	r.UpdatedBy = userID

	return
}

func GetUserID(tx *gorm.DB) uint {
	userID, ok := tx.Statement.Context.Value("user_id").(uint)
	if !ok {
		return 0
	}
	return userID
}

