package rab

import (
	"context"
)

type rabService struct {
	rabRepo Repository
}

func NewRabService(repo Repository) Service {
	return &rabService{
		rabRepo: repo,
	}
}

func (s *rabService) GetListRab(ctx context.Context, q string, code string) ([]Rab, error) {
	return s.rabRepo.FetchAll(ctx, q, code)
}
