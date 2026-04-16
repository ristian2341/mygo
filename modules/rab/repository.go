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

func (r *mysqlRabRepository) Create(ctx context.Context, rab *Rab) error {
	return r.db.WithContext(ctx).Create(rab).Error
}

func (r *mysqlRabRepository) GetAll(ctx context.Context) ([]Rab, error) {
	var rabs []Rab
	err := r.db.WithContext(ctx).Find(&rabs).Error
	return rabs, err
}
