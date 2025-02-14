package usecase

import (
	"context"

	"go.uber.org/zap"
)

type Repository interface {
	Get(context.Context, string) (string, error)
	Create(context.Context, string, string) error
	LenRows(context.Context) (int, error)
}

type Usecase struct {
	log     *zap.Logger
	Repo    Repository
}

func NewUsecase(log *zap.Logger, Repo Repository) (*Usecase, error) {
	log.Named("usecase")

	return &Usecase{
		log:     log,
		Repo:    Repo,
	}, nil
}