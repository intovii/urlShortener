package repository

import (
	"URLShortener/config"
	inMemory "URLShortener/internal/repository/inMemo"
	"URLShortener/internal/repository/postgres"
	"context"
	"errors"

	"go.uber.org/zap"
)

type Repository interface {
	Get(context.Context, string) (string, error)
	Create(context.Context, string, string) error
	OnStart(context.Context) error
	OnStop(context.Context) error
	LenRows(context.Context) (int, error)
}


func NewRepository(ctx context.Context, cfg *config.ConfigModel, log *zap.Logger) (Repository, error) {
	log.Named("repository")

	var repo Repository
	switch cfg.StorageType {
	case "PSQL":
		log.Info("PSQL storage")
		repo, _ = postgres.NewPostgresRepository(ctx, cfg, log)
	case "InMemo":
		log.Info("In-memory storage")
		repo, _ = inMemory.NewInMemoRepository(ctx, log)
	default:
		log.Error("Storage type invalid")
		return nil, errors.New("invalid storage type")
	}
	
	return repo, nil
}