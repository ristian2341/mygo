package models

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