package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	httpgin "github.com/MishaZhem/Gona/src/gateway/internal/ports/http"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	userserviceAddr := getEnv("USERSERVICE_ADDR", "localhost:50051")
	httpPort := getEnv("GATEWAY_HTTP_PORT", "8080")

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

	conn, err := grpc.Dial(userserviceAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to userservice: %v", err)
	}
	defer conn.Close()

	app := app.NewApp(conn, 5*time.Second)

	httpServer := httpgin.NewHTTPServer(httpPort, app)

	eg.Go(func() error {
		fmt.Println("starting HTTP server")
		errCh := make(chan error)

		defer func() {
			fmt.Println("stopping HTTP server")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("error on HTTP server closing occurred: %s", err.Error())
			}
			close(errCh)
		}()

		go func() {
			if err := httpServer.Listen(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("HTTP server error: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("servers shutdown: %s\n", err.Error())
	}
}

func getEnv(key, init string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return init
}
