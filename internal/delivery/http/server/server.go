package server

import (
	"URLShortener/domain"
	protos "URLShortener/pb"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (s *Server) OnStart(_ context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Server.Host+":"+s.cfg.Server.GRPCPort)
	if err != nil {
		s.log.Error("failed to listen: ", zap.Error(err))
		return fmt.Errorf("failed to listen:  %w", err)
	}

	protos.RegisterGatewayServer(s.RPC, s)

	reflection.Register(s.RPC)
	go func() {
		s.log.Debug("grpc serv started")
		if err = s.RPC.Serve(lis); err != nil {
			s.log.Error("failed to serve: " + err.Error())
		}
		return
	}()
	
	conn, err := grpc.DialContext(
		context.Background(),
		s.cfg.Server.Host + ":" + s.cfg.Server.GRPCPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.log.Error("Failed to dial server:", zap.Error(err))
		return err
	}
	defer conn.Close()
	
	gwmux := runtime.NewServeMux()
	err = protos.RegisterGatewayHandler(context.Background(), gwmux, conn)
	if err != nil {
		s.log.Error("Failed to register gateway:", zap.Error(err))
		return err
	}
	
	s.log.Info(fmt.Sprintf("Serving gRPC-Gateway on port %s", s.cfg.Server.HTTPPort))
	gwServer := &http.Server{
		Addr:    ":" + s.cfg.Server.HTTPPort,
		Handler: gwmux,
	}
	
	
	if err = gwServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.log.Error("Server closed:", zap.Error(err))
			return err
		}	
		s.log.Error("Failed to listen and serve:", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.log.Info("stop grpc server")
	s.RPC.GracefulStop()
	return nil
}

func (s *Server) Create(ctx context.Context, request *protos.CreateUrlRequest) (*protos.CreateUrlResponse, error) {
	shortUrl, err := s.Usecase.Create(ctx, request.Url)
	if err != nil {
		s.log.Error("failed to create", zap.Error(err))
		if errors.Is(err, domain.ErrInvalidArgument) {
			return nil, status.Error(http.StatusBadRequest, "Bad request")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &protos.CreateUrlResponse{ShortUrl: shortUrl}, nil
}

func (s *Server) Get(ctx context.Context, request *protos.GetUrlRequest) (*protos.GetUrlResponse, error) {
	originalUrl, err := s.Usecase.Get(ctx, request.Url)
	if err != nil {
		s.log.Error("failed to get original url:", zap.Error(err))
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(http.StatusBadRequest, "Bad request")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &protos.GetUrlResponse{OriginalUrl: originalUrl}, nil
}