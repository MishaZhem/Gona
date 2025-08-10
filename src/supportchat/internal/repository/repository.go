package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MishaZhem/Gona/src/supportchat/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRepository(client *redis.Client, ttl time.Duration) *Repository {
	return &Repository{
		client: client,
		ttl:    ttl,
	}
}

func (r *Repository) SaveMessage(ctx context.Context, userID string, msg *domain.Message) error {
	key := fmt.Sprintf("chat:%s:messages", userID)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()
	pipe.RPush(ctx, key, data)
	pipe.Expire(ctx, key, r.ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *Repository) GetMessages(ctx context.Context, userID string) ([]*domain.Message, error) {
	key := fmt.Sprintf("chat:%s:messages", userID)
	items, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var msgs []*domain.Message
	for _, item := range items {
		var msg domain.Message
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			continue
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func (r *Repository) SaveChatMeta(ctx context.Context, meta *domain.ChatMeta) error {
	key := fmt.Sprintf("chat:%s:meta", meta.UserID)
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *Repository) GetChatMeta(ctx context.Context, userID string) (*domain.ChatMeta, error) {
	key := fmt.Sprintf("chat:%s:meta", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var meta domain.ChatMeta
	if err := json.Unmarshal([]byte(val), &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (r *Repository) DeleteChat(ctx context.Context, userID string) error {
	keyMsgs := fmt.Sprintf("chat:%s:messages", userID)
	keyMeta := fmt.Sprintf("chat:%s:meta", userID)
	return r.client.Del(ctx, keyMsgs, keyMeta).Err()
}

func (r *Repository) GetActiveChats(ctx context.Context) ([]*domain.ChatMeta, error) {
	var chats []*domain.ChatMeta
	iter := r.client.Scan(ctx, 0, "chat:*:meta", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var meta domain.ChatMeta
		if err := json.Unmarshal([]byte(val), &meta); err != nil {
			continue
		}
		if meta.WorkerID == "" && meta.NeedsWorker {
			chats = append(chats, &meta)
		}
		chats = append(chats, &meta)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error scanning active chats: %w", err)
	}

	return chats, nil
}
