package domain

import (
	"time"
)

type Auth struct {
	Email string
	Password string
	Role string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
