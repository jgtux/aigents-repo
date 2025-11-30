package services

import (
	a_at "aigents-base/internal/auth-land/auth/atoms"
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"
	c_at "aigents-base/internal/common/atoms"


	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)


type AuthService struct {
	r auitf.AuthRepositoryITF
}

func NewAuthService(repo auitf.AuthRepositoryITF) auitf.AuthServiceITF {
	return &AuthService{r: repo}
}

func (s *AuthService) Create(gctx *gin.Context, data *d.Auth) error {
	hashedPass := a_at.HashPassAtom(data.Password)
	data.Password = hashedPass

	if err := s.r.Create(gctx, data); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Comparate(gctx *gin.Context, data *d.Auth) error {
	auth := &d.Auth{}
	auth.Email = data.Email

	err := s.r.GetByEmail(gctx, auth)
	if err != nil {
		if err.Error() == "ERR_EMAIL_NOT_FOUND" {
			err = c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusUnauthorized,
				"(S) Invalid credentials.",
				fmt.Sprintf("Authentication of %s not found", data.Email))
		}
		return err
	}

	hashedTriedPass := a_at.HashPassAtom(data.Password)

	if hashedTriedPass != auth.Password {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusUnauthorized,
			"(S) Invalid credentials.",
			fmt.Sprintf("Incorrect password of %s", auth.Email))
		return err
	}

	return nil
}

func (s *AuthService) GetByID(gctx *gin.Context, data *d.Auth) error {
	return s.r.GetByID(gctx, data)
}

func (s *AuthService) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Auth, error) {
	return s.r.Fetch(gctx, limit, offset)
}

func (s *AuthService) Update(gctx *gin.Context, data *d.Auth) error {
	return s.r.Update(gctx, data)
}

func (s *AuthService) Delete(gctx *gin.Context, data *d.Auth) error {
	return s.r.Delete(gctx, data)
}
