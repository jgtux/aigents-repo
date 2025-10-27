package domain

import (
	"time"
)

type User struct {
	UUID string
	AuthUUID string
	FirstName string
	LastName string
	document_id string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
