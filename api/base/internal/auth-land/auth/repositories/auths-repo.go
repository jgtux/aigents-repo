package repositories

import (
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"

	"github.com/gin-gonic/gin"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auitf.AuthRepositoryITF {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(data *d.Auth) func() {
	return nil
}

func (a *AuthRepository) GetByEmail(data *d.Auth) func() {
	return nil
}

func (a *AuthRepository) GetByID(data *d.Auth) func() {
	return nil
}
func (a *AuthRepository) Fetch(limit, offset int) ([]d.Auth, func()) {
	return []d.Auth{}, nil
}

func (a *AuthRepository) Update(data *d.Auth) func() {
	return nil
}

func (a *AuthRepository) Delete(data *d.Auth) func() {
	return nil
}
