package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreatedAt time.Time
	ID        uuid.UUID
	Email     string
	Password  string
}
