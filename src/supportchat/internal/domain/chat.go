package domain

import (
	"time"

	uuid "github.com/google/uuid"
)

type Message struct {
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatMeta struct {
	UserID      string    `json:"user_id"`
	WorkerID    string    `json:"worker_id,omitempty"`
	NeedsWorker bool      `json:"needs_worker"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GenerateID() string {
	return uuid.New().String()
}
