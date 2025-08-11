package domain

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"time"`
	AvatarURL string    `json:"avatar_url"`
	Worker    bool      `json:"worker"`
}
