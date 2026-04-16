package rab

import (
	"context"

	"gorm.io/gorm"
)

type mysqlRabRepository struct {
	db *gorm.DB
}

func NewMySQLRabRepository(db *gorm.DB) Repository {
	return &mysqlRabRepository{
		db: db,
	}
}

func (r *mysqlRabRepository) FetchAll(ctx context.Context, q string, code string) ([]Rab, error) {
	var dataRabs []Rab
	query := r.db.WithContext(ctx).Model(&dataRabs)

	if q != "" {
		searchTerm := "%" + q + "%"
		query = query.Where("CodeProject LIKE ? OR email LIKE ? OR nama LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if code != "" {
		query = query.Where("code = ?", code)
	}

	err := query.Find(&dataRabs).Error
	return dataRabs, err
}
