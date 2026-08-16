package model

import (
	"time"
	"github.com/google/uuid"
)

type ChatRole string
const (
	CHAT_ROLE_AI   ChatRole = "ai"
	CHAT_ROLE_USER ChatRole = "user"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	Body   string `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	SenderRole ChatRole `json:"sender_role"`
}