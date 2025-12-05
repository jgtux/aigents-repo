package services

import (
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"


	"github.com/gin-gonic/gin"
	"fmt"
)


type AgentService struct {
	r agitf.AgentRepositoryITF
}

func NewAgentService(repo agitf.AgentRepositoryITF) agitf.AgentServiceITF {
	return &AgentService{r: repo}
}

func (s *AgentService) Create(gctx *gin.Context, data *d.Agent) error {
	data.AgentConfig.AgentSystem.SystemPreset = map[string]any{
		"system_prompt": fmt.Sprintf("You're a helpful assistant and your job will be doing this description: %s", data.Description),
	}

	return s.r.Create(gctx, data)
}

func (s *AgentService) FetchAgentsByLoggedAuth(gctx *gin.Context, authUUID string, limit, offset uint64) ([]d.Agent, error) {
	return s.r.FetchAgentsByLoggedAuth(gctx, authUUID, limit, offset)
}

func (s *AgentService) FetchCategories(gctx *gin.Context) ([]d.AgentCategory, error) {
	return s.r.FetchCategories(gctx)
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

func (a *AgentService) FetchWithFilter(gctx *gin.Context, flags []string, limit, offset uint64) ([]d.Agent, error) {
	return []d.Agent{}, nil
}
