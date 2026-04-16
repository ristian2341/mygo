package rab

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Rab struct {
	Code        string         `gorm:"primaryKey;type:char(12)" json:"code"`
	CodeProject string         `json:"code_project"`
	TglMulai    time.Time      `json:"tgl_mulai"`
	TglAkhir    time.Time      `json:"tgl_akhir"`

	CreatedBy uint `json:"created_by"`
	UpdatedBy uint `json:"updated_by"`

	CreatedAt time.Time
	UpdatedAt time.Time

	RabDetails []RabDetail `gorm:"foreignKey:Code;references:Code"`
	RabKelompoks []RabKelompok `gorm:"foreignKey:Code;references:Code"`
}

func (Rab) TableName() string {
	return "rab"
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

type RabDetail struct {
	Code         string  `gorm:"primaryKey;type:char(12)" json:"code"`
	Line         uint    `gorm:"primaryKey;type:int(11)" json:"line"`
	CodeKelompok string  `gorm:"primaryKey;type:char(12)" json:"code_kelompok"`
	Uraian       string  `gorm:"type:varchar(250)" json:"uraian"`
	Volume       float64 `gorm:"type:double" json:"volume"`
	Satuan       string  `gorm:"type:varchar(100)" json:"satuan"`
	HargaSatuan  float64 `gorm:"type:double" json:"harga_satuan"`
	Keterangan   string  `gorm:"type:varchar(200)" json:"keterangan"`
}

func (RabDetail) TableName() string {
	return "rab_detail"
}

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

func GetUserID(tx *gorm.DB) uint {
	userID, ok := tx.Statement.Context.Value("user_id").(uint)
	if !ok {
		return 0
	}
	return userID
}

type Repository interface {
	FetchAll(ctx context.Context, q string, code string) ([]Rab, error)
}

type Service interface {
	GetListRab(ctx context.Context, q string, code string) ([]Rab, error)
}
