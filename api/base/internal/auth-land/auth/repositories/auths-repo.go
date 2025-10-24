package repositories

import (
	c_at "aigents-base/internal/common/atoms"
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"database/sql"
	"net/http"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auitf.AuthRepositoryITF {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(data *d.Auth) func(*gin.Context) {
	query := `
		INSERT INTO auths (
			email,
			password
		) VALUES ($1, $2)
		RETURNING auth_uuid, created_at, updated_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00');
	`

	err := a.db.QueryRow(query, data.Email, data.Password).Scan(
		&data.UUID,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return c_at.RespFuncAbortAtom(
					http.StatusConflict,
					"(R) Email already registered.",
				)
			}
		}

		return c_at.RespFuncAbortAtom(
			http.StatusInternalServerError,
			"(R) An unknown error occurred.",
		)
	}

	return nil
}

func (a *AuthRepository) GetByEmail(data *d.Auth) func(*gin.Context) {
	query := `SELECT auth_uuid,
                         password,
                         created_at,
                         updated_at,
                         deleted_at
                  FROM auths
                  WHERE email = $1;`

	err := a.db.QueryRow(query, data.Email).Scan(
		&data.UUID,
		&data.Password,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c_at.RespFuncAbortAtom(
				http.StatusNotFound,
				"(R) Authentication not found.")
		} else {
			return c_at.RespFuncAbortAtom(
				http.StatusInternalServerError,
				"(R) An unknown error ocurred.")
		}
	}

	return nil
}

func (a *AuthRepository) GetByID(data *d.Auth) func(*gin.Context) {

	return nil
}
func (a *AuthRepository) Fetch(limit, offset int) ([]d.Auth, func(*gin.Context)) {
	return []d.Auth{}, nil
}

func (a *AuthRepository) Update(data *d.Auth) func(*gin.Context) {
	return nil
}

func (a *AuthRepository) Delete(data *d.Auth) func(*gin.Context) {
	return nil
}
