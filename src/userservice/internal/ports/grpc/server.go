package grpc

import (
	"userservice/internal/app"
	"userservice/internal/middleware"

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
