package repositories

import (
	d "aigents-base/internal/users/domain"
	uitf "aigents-base/internal/users/interfaces"
	c_at "aigents-base/internal/common/atoms"

	"database/sql"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UsersRepository struct {
	db *sql.DB
}


func NewUsersRepository(db *sql.DB) uitf.UsersRepositoryITF {
	return &UsersRepository{db: db}
}

func (a *UsersRepository) Create(data *d.User) func(*gin.Context) {
	query := `
		INSERT INTO users (
			auth_uuid,
			first_name,
                        last_name,
                        document_id,
		) VALUES ($1, $2, $3, $4)
		RETURNING user_uuid, created_at, updated_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00');
	`

	err := a.db.QueryRow(query,
		data.AuthUUID,
		data.FirstName,
		data.LastName,
		data.DocumentID,
	).Scan(
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
					"(R) Authentication already linked to another user.",
				)
			}
		}

		return c_at.RespFuncAbortAtom(
			http.StatusInternalServerError,
			fmt.Sprintf("(R) An unknown error occurred: %s", err.Error()),
		)
	}

	return nil
}


func (a *UsersRepository) GetByID(data *d.User) func(*gin.Context) {

	return nil
}
func (a *UsersRepository) Fetch(limit, offset int) ([]d.User, func(*gin.Context)) {
	return []d.User{}, nil
}

func (a *UsersRepository) Update(data *d.User) func(*gin.Context) {
	return nil
}

func (a *UsersRepository) Delete(data *d.User) func(*gin.Context) {
	return nil
}
