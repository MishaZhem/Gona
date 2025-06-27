package middleware

import (
	"context"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(logger log.FieldLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(log.Fields{
					"panic":  r,
					"stack":  string(debug.Stack()),
					"method": info.FullMethod,
				}).Error("panic recovered in gRPC")
			}
		}()

		return handler(ctx, req)
	}
}
