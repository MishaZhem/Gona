package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MishaZhem/Gona/src/userservice/internal/adapters/postgres"
	"github.com/MishaZhem/Gona/src/userservice/internal/app"
	grpcPort "github.com/MishaZhem/Gona/src/userservice/internal/ports/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	port := getEnv("GRPC_PORT", "50051")
	postgresURL := getEnv("POSTGRES_URL", "postgres://postgres:postgres@db:5432/userdb?sslmode=disable")
	secretKey := getEnv("JWT_SECRET", "SomeHardAndSecretWord")
	tokenTTL := time.Hour * 24

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})

	pool, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer pool.Close()

	repo := postgres.NewUserRepository(pool, logger)
	tokenService := app.NewJWTService(secretKey, tokenTTL)
	app := app.NewApp(repo, tokenService, logger)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.WithError(err).Fatal("can't create listener")
		return
	}

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	grpcServer := grpcPort.NewGRPCServer(app, logger)

	eg.Go(func() error {
		logger.Infof("starting gRPC server on port %s", port)
		errCh := make(chan error)

		defer func() {
			log.Infof("stopping GRPC server")
			grpcServer.GracefulStop()
			_ = lis.Close()
			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			logger.WithError(err).Error("gRPC server failed")
			return err
		}
	})

	if err := eg.Wait(); err != nil {
		logger.WithError(err).Fatal("server shutdown with error")
	}
}

func getEnv(key, init string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return init
}
