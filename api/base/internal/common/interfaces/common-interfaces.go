package interfaces

import (
	"github.com/gin-gonic/gin"
)

type Common[T any] interface {
	Create(gctx *gin.Context, data *T) error
	GetByID(gctx *gin.Context, data *T) error
	Fetch(gctx *gin.Context, limit, offset uint64) ([]T, error)
	Update(gctx *gin.Context, data *T) error
	Delete(gctx *gin.Context, data *T) error
}
