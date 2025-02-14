package app

import (
	"URLShortener/config"
	"log/slog"

	"URLShortener/internal/delivery/http/server"
	"URLShortener/internal/repository"
	"URLShortener/internal/usecase"

	"context"

	"go.uber.org/zap"
)


type App struct {
	server		*server.Server
	usecase    	*usecase.Usecase
	repository 	repository.Repository
	ctx context.Context
	cfg *config.ConfigModel
}

func New() *App {	
	ctx := context.Background()
	
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("error with config")
		return nil
	}
	
	log, err := zap.NewProduction()
	if err != nil {
		slog.Error("error with zap logger")
		return nil
	}
	
	repo, _ := repository.NewRepository(ctx, cfg, log)

	usecase, _ := usecase.NewUsecase(log, repo)

	server, _ := server.NewServer(ctx, cfg, log, usecase)

	return &App{
		ctx: ctx,
		repository: repo,
		usecase: usecase,
		server: server,
	}
}

func (a *App) Run() {
	a.repository.OnStart(a.ctx)
	defer a.repository.OnStop(a.ctx)

	a.usecase.OnStart(a.ctx)
	defer a.usecase.OnStop(a.ctx)

	a.server.OnStart(a.ctx)
	defer a.server.OnStop(a.ctx)
}