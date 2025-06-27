package grpc

import (
	"user-service/internal/app"
	"user-service/internal/middleware"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

func NewGRPCServer(app app.App, logger *log.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(logger),
			middleware.LoggingInterceptor(logger),
		),
	)
	Register(server, app)
	return server
}
