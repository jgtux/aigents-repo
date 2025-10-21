package interfaces

import (
	"github.com/gin-gonic/gin"
	"context"
)

type Common[T any, PK ~string | ~uint64] interface {
	Create(data *T) error
	GetByID(PK) (*T, error)
	Fetch(limit, offset int) ([]T, error)
	Update(data *T) error
	Delete(id PK) error
}
