package repositories

import (
	c_at "aigents-base/internal/common/atoms"
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"

	"github.com/lib/pq"
	"strings"
	"database/sql"
	"net/http"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auitf.AuthRepositoryITF {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(data *d.Auth) func() {
	query := `
		INSERT INTO auth (
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
				c_at.RespFuncAbortAtom(
					http.StatusConflict,
					"(R) Email already registered.",
				)
				return nil
			}
		}

		c_at.RespFuncAbortAtom(
			http.StatusInternalServerError,
			"(R) An unknown error occurred.",
		)
		return nil
	}

	return nil
}

func (a *AuthRepository) GetByEmail(data *d.Auth) func() {
	query := `SELECT auth_uuid,
                         password,
                         created_at,
                         updated_at,
                         deleted_at
                  FROM auth
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
			c_at.RespFuncAbortAtom(
				http.StatusNotFound,
				"(R) Authentication not found.")
		} else {
			c_at.RespFuncAbortAtom(
				http.StatusInternalServerError,
				"(R) An unknown error ocurred.")
		}
	}

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
