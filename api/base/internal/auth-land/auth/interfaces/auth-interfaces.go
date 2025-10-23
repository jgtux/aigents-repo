package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/auth-land/auth/domain"

	"github.com/gin-gonic/gin"
)

type AuthServiceITF interface {
	citf.Common[d.Auth]
	Comparate(data *d.Auth) func(*gin.Context)
}

type AuthRepositoryITF interface {
	citf.Common[d.Auth]
	GetByEmail(data *d.Auth) func(*gin.Context)
}
