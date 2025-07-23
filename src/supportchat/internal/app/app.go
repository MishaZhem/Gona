package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MishaZhem/Gona/src/supportchat/internal/domain"
	"github.com/MishaZhem/Gona/src/supportchat/internal/repository"
)

type App struct {
	repo *repository.Repository
	bot  SupportBot
}

type SupportBot interface {
	GenerateResponse(input string) (string, bool)
}

func NewApp(repo *repository.Repository, bot SupportBot) *App {
	return &App{
		repo: repo,
		bot:  bot,
	}
}

func (a *App) StartNewChat(ctx context.Context, userID string) (*domain.ChatMeta, error) {
	if err := a.repo.DeleteChat(ctx, userID); err != nil {
		return nil, err
	}

	meta := &domain.ChatMeta{
		UserID:      userID,
		NeedsWorker: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return meta, a.repo.SaveChatMeta(ctx, meta)
}

func (a *App) SendUserMessage(ctx context.Context, userID, text string) ([]*domain.Message, error) {
	meta, err := a.repo.GetChatMeta(ctx, userID)
	if err != nil {
		if meta, err = a.StartNewChat(ctx, userID); err != nil {
			return nil, err
		}
	}

	msg := &domain.Message{
		From:      "user",
		Text:      text,
		Timestamp: time.Now(),
	}
	if err := a.repo.SaveMessage(ctx, userID, msg); err != nil {
		return nil, err
	}

	fmt.Print(meta.WorkerID)
	fmt.Print(meta.WorkerID == "")
	if meta.WorkerID == "" {
		answer, ok := a.bot.GenerateResponse(text)
		botMsg := &domain.Message{
			From:      "bot",
			Text:      answer,
			Timestamp: time.Now(),
		}
		if err := a.repo.SaveMessage(ctx, userID, botMsg); err != nil {
			return nil, err
		}
		meta.NeedsWorker = ok
	}

	meta.UpdatedAt = time.Now()
	_ = a.repo.SaveChatMeta(ctx, meta)

	return a.repo.GetMessages(ctx, userID)
}

func (a *App) SendWorkerMessage(ctx context.Context, userID, workerID, text string) ([]*domain.Message, error) {
	meta, err := a.repo.GetChatMeta(ctx, userID)
	if err != nil {
		if meta, err = a.StartNewChat(ctx, userID); err != nil {
			return nil, err
		}
	}

	if meta.WorkerID != "" && meta.WorkerID != workerID {
		return nil, errors.New("chat is already handled by another worker")
	}

	meta.WorkerID = workerID
	meta.NeedsWorker = false
	meta.UpdatedAt = time.Now()
	if err := a.repo.SaveChatMeta(ctx, meta); err != nil {
		return nil, err
	}

	msg := &domain.Message{
		From:      workerID,
		Text:      text,
		Timestamp: time.Now(),
	}
	if err := a.repo.SaveMessage(ctx, userID, msg); err != nil {
		return nil, err
	}

	return a.repo.GetMessages(ctx, userID)
}

func (a *App) GetActiveChats(ctx context.Context) ([]*domain.ChatMeta, error) {
	return a.repo.GetActiveChats(ctx)
}

func (a *App) GetMessages(ctx context.Context, userID string) ([]*domain.Message, error) {
	return a.repo.GetMessages(ctx, userID)
}

func (a *App) DeleteChat(ctx context.Context, userID string) error {
	return a.repo.DeleteChat(ctx, userID)
}

func (a *App) GetChatMeta(ctx context.Context, userID string) (*domain.ChatMeta, error) {
	return a.repo.GetChatMeta(ctx, userID)
}
