package models

import (
	"time"

	"gorm.io/gorm"
)

type RabKelompok struct {
	Code     string `gorm:"primaryKey;type:char(12)" json:"code"`
	Kelompok string `json:"kelompok"`

	CreatedBy uint `json:"created_by"`
	UpdatedBy uint `json:"updated_by"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RabKelompok) TableName() string {
	return "rab_kelompok"
}

func (r *RabKelompok) BeforeCreate(tx *gorm.DB) (err error) {
	userID := GetUserID(tx)

	r.CreatedBy = userID
	r.UpdatedBy = userID

	return
}

func (r *RabKelompok) BeforeUpdate(tx *gorm.DB) (err error) {
	userID := GetUserID(tx)

	r.UpdatedBy = userID

	return
}
