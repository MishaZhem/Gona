package middleware

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger log.FieldLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		logger.WithFields(log.Fields{
			"method":  info.FullMethod,
			"latency": time.Since(start),
			"error":   err,
		}).Info("gRPC request")

		return resp, err
	}
}
