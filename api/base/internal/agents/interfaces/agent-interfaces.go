package interfaces

import (
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/agents/domain"

	"github.com/gin-gonic/gin"
)

type AgentServiceITF interface {
	citf.Common[d.Agent]
	FetchWithFilter(gctx *gin.Context, flags []string, limit, offset uint64) ([]d.Agent, error)
}

type AgentRepositoryITF interface {
	citf.Common[d.Agent]
	FetchWithFilter(gctx *gin.Context, flags []string, limit, offset uint64) ([]d.Agent, error)
	GetAgentByUUID(gctx *gin.Context, agentUUID string) (*d.Agent, error)
}
