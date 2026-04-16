package rab

import "context"

type Rab struct {
	ID          uint   `json:"id"`
	ProjectName string `json:"project_name"`
	TotalAmount int64  `json:"total_amount"`
	CreatedBy   string `json:"created_by"` // Username dari modul user
}

func (Rab) TableName() string {
	return "rabs"
}

type Repository interface {
	Create(ctx context.Context, rab *Rab) error
	GetAll(ctx context.Context) ([]Rab, error)
}

type Service interface {
	CreateRab(ctx context.Context, projectName string, totalAmount int64, username string) error
	GetListRab(ctx context.Context) ([]Rab, error)
}
