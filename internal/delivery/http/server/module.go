package server

import (
	"URLShortener/config"
	protos "URLShortener/pb"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)


type Usecase interface{
	Create(context.Context, string) (string, error)
	Get(context.Context, string) (string, error)
}

type Server struct {
	ctx 	context.Context
	log  	*zap.Logger
	cfg     *config.ConfigModel
	RPC     *grpc.Server
	Usecase Usecase
	protos.UnimplementedGatewayServer
}

func NewServer(ctx context.Context, cfg *config.ConfigModel, log *zap.Logger, uc Usecase) (*Server, error) {
	log.Named("server")

	return &Server{
		ctx: 		ctx,
		log:  		log,
		cfg:    	cfg,
		RPC:    	grpc.NewServer(),
		Usecase: 	uc,
	}, nil
}