package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/auth-land/auth/domain"
)

type AuthServiceITF interface {
	citf.Common[d.Auth, string]
	Login(*d.Auth) error
}

type AuthRepositoryITF interface {
	citf.Common[d.Auth, string]
	GetByEmail(*d.Auth) error
}
