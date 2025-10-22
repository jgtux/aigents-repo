package services

import (
	a_at "aigents-base/internal/auth-land/auth/atoms"
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"
	c_at "aigents-base/internal/common/atoms"

	"net/http"
)

type AuthService struct {
	r auitf.AuthRepositoryITF
}

func NewAuthService(repo auitf.AuthRepositoryITF) auitf.AuthServiceITF {
	return &AuthService{r: repo}
}

func (s *AuthService) Create(data *d.Auth) func() {
	hashedPass := a_at.HashPassAtom(data.Password)
	data.Password = hashedPass

	return s.r.Create(data)
}

func (s *AuthService) Comparate(data *d.Auth) func() {
	auth := &d.Auth{}
	auth.Email = data.Email

	err := s.r.GetByEmail(auth)
	if err != nil {
		return err
	}

	hashedTriedPass := a_at.HashPassAtom(data.Password)

	if hashedTriedPass != auth.Password {
		return c_at.RespFuncAbortAtom(
			http.StatusUnauthorized,
			"(S) Invalid credentials.")
	}

	return nil
}

func (s *AuthService) GetByID(data *d.Auth) func() {
	return s.r.GetByID(data)
}

func (s *AuthService) Fetch(limit, offset int) ([]d.Auth, func()) {
	return s.r.Fetch(limit, offset)
}

func (s *AuthService) Update(data *d.Auth) func() {
	return s.r.Update(data)
}

func (s *AuthService) Delete(data *d.Auth) func() {
	return s.r.Delete(data)
}
