package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/agents/domain"

	"github.com/gin-gonic/gin"
)

type AgentServiceITF interface {
	citf.Common[d.Agent]
	Search(gctx *gin.Context, flags []string) ([]d.Agent, error)
}

type AgentRepositoryITF interface {
	citf.Common[d.Agent]
	FetchWithFilter(gctx *gin.Context, flags []string) ([]d.Agent, error)
}
