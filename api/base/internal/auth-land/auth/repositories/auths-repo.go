package repositories

import (
	d "aigents-base/internal/auth-land/auth/domain"
	itf "aigents-base/internal/common/interfaces"


	"github.com/gin-gonic/gin"
	"database/sql"
	"context"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) itf.Common[d.Auth, string] {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(gctx *gin.Context, data *d.Auth) error {
	return nil
}

func (a *AuthRepository) GetByID(gctx *gin.Context, uuid string) (*d.Auth, error) {
	return &d.Auth{}, nil
}
func (a *AuthRepository) Fetch(gctx *gin.Context, limit, offset int) ([]d.Auth, error) {
	return []d.Auth{}, nil
}

func (a *AuthRepository) Update(gctx *gin.Context, data *d.Auth) error {
	return nil
}

func (a *AuthRepository) Delete(gctx *gin.Context, uuid string) error {
	return nil
}
