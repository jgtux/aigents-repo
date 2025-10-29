package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/users/domain"

	"github.com/gin-gonic/gin"
)


type UsersServiceITF interface {
	citf.Common[d.User]
}

type UsersRepositoryITF interface {
	citf.Common[d.User]
}
