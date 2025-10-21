package services

import (
	d "aigents-base/internal/auth-land/auth/domain"
	itf "aigents-base/internal/common/interfaces"

	"github.com/gin-gonic/gin"
	"context"
)

type AuthService struct {
	r itf.Common[d.Auth, string]
}

func NewAuthService(repo itf.Common[d.Auth, string]) itf.Common[d.Auth, string] {
	return &AuthService{r: repo}
}

func (s *AuthService) Create(data *d.Auth) error {
	return s.r.Create(data)
}

func (s *AuthService) Login(data *d.Auth) error {
	err := s.r.GetByEmail(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetByID(id string) (*d.Auth, error) {
	return s.r.GetByID(id)
}

func (s *AuthService) Fetch(limit, offset int) ([]d.Auth, error) {
	return s.r.Fetch(limit, offset)
}

func (s *AuthService) Update(data *d.Auth) error {
	return s.r.Update(data)
}

func (s *AuthService) Delete(id string) error {
	return s.r.Delete(id)
}
