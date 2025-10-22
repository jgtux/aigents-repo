package domain

import (
	"time"
)

type Auth struct {
	UUID string
	Email string
	Password string
	Role string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
