package rab

import (
	"context"
	"errors"
)

type rabService struct {
	rabRepo Repository
}

func NewRabService(repo Repository) Service {
	return &rabService{
		rabRepo: repo,
	}
}

func (s *rabService) CreateRab(ctx context.Context, projectName string, totalAmount int64, username string) error {
	if projectName == "" {
		return errors.New("Project name is required")
	}

	rab := &Rab{
		ProjectName: projectName,
		TotalAmount: totalAmount,
		CreatedBy:   username,
	}

	return s.rabRepo.Create(ctx, rab)
}

func (s *rabService) GetListRab(ctx context.Context) ([]Rab, error) {
	return s.rabRepo.GetAll(ctx)
}
