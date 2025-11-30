package services

import (
	a_at "aigents-base/internal/auth-land/auth/atoms"
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"
	c_at "aigents-base/internal/common/atoms"


	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)


type AgentService struct {
	r agitf.AgentRepositoryITF
}

func NewAuthService(repo agitf.AgentRepositoryITF) agitf.AgentServiceITF {
	return &AgentService{r: repo}
}

func (s *AgentService) Create(gctx *gin.Context, data *d.Agent) error {
	return nil
}

func (a *AgentService) FetchWithFilter(gctx *gin.Context, flags []string, limit, offset uint64) ([]d.Agent, error) {
	return []d.Agent{}, nil
}

func (s *AgentService) GetByID(gctx *gin.Context, data *d.Agent) error {
	return s.r.GetByID(gctx, data)
}

func (s *AgentService) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Agent, error) {
	return s.r.Fetch(gctx, limit, offset)
}

func (s *AgentService) Update(gctx *gin.Context, data *d.Agent) error {
	return s.r.Update(gctx, data)
}

func (s *AgentService) Delete(gctx *gin.Context, data *d.Agent) error {
	return s.r.Delete(gctx, data)
}
