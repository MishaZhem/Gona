package main

import (
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/MishaZhem/Gona/src/supportchat/internal/app"
	"github.com/MishaZhem/Gona/src/supportchat/internal/ports/ws"
	"github.com/MishaZhem/Gona/src/supportchat/internal/repository"
	"github.com/redis/go-redis/v9"
)

type DummyBot struct{}

func (b *DummyBot) GenerateResponse(input string) (string, bool) {
	if input == "Ask" {
		return "I'll get a support operator on the line now", false
	}
	return "Bot: I got your message!", true
}

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: "", // если есть пароль, укажи
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})

	repo := repository.NewRepository(rdb, 24*time.Hour)
	bot := &DummyBot{}
	appService := app.NewApp(repo, bot, logger)
	wsHandler := ws.NewHandler(appService)

	http.HandleFunc("/chat", wsHandler.ServeWS)

	port := getEnv("SUPPORTCHAT_PORT", "9000")
	log.Printf("Starting supportchat WebSocket server on :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
