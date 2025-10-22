package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/auth-land/auth/domain"
)

type AuthServiceITF interface {
	citf.Common[d.Auth]
	Comparate(data *d.Auth) func()
}

type AuthRepositoryITF interface {
	citf.Common[d.Auth]
	GetByEmail(data *d.Auth) func()
}
