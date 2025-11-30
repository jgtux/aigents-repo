package repositories

import (
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"
	c_at "aigents-base/internal/common/atoms"
	"fmt"

	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auitf.AuthRepositoryITF {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(gctx *gin.Context, data *d.Auth) error {
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
				err = c_at.AbortAndBuildErrLogAtom(
					gctx,
					http.StatusConflict,
					"(R) Email already registered.",
					fmt.Sprintf("Email %s already registred", data.Email))
				return err
			}
		}


		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(R) Could not register authentication.",
			fmt.Sprintf("An unknown error ocurred: %s", err.Error()))
		return err
	}

	return nil
}

func (a *AuthRepository) GetByEmail(gctx *gin.Context, data *d.Auth) error {
	query := `SELECT auth_uuid,
                         password,
                         created_at,
                         updated_at,
        	         COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00')
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
			err = c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusUnauthorized,
				"(R) Authentication not found.",
				fmt.Sprintf("Authentication of %s not found", data.Email))
			return err
		}


		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(R) Could not find authentication.",
			fmt.Sprintf("An unknown error ocurred: %s", err.Error()))
		return err
	}

	return nil
}

func (a *AuthRepository) GetByID(gctx *gin.Context, data *d.Auth) error {

	return nil
}
func (a *AuthRepository) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Auth, error) {
	return []d.Auth{}, nil
}

func (a *AuthRepository) Update(gctx *gin.Context, data *d.Auth) error {
	return nil
}

func (a *AuthRepository) Delete(gctx *gin.Context, data *d.Auth) error {
	return nil
}
