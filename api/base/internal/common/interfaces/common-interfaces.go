package interfaces

import (
	"github.com/gin-gonic/gin"
)

type Common[T any] interface {
	Create(data *T) func(*gin.Context)
	GetByID(data *T) func(*gin.Context)
	Fetch(limit, offset int) ([]T, func(*gin.Context))
	Update(data *T) func(*gin.Context)
	Delete(data *T) func(*gin.Context)
}
