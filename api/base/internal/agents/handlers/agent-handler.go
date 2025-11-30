
package handlers

import (
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"
	c_at "aigents-base/internal/common/atoms"

	"net/http"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	s agitf.AgentServiceITF
}

func NewAuthHandler(sv agitf.AgentServiceITF) *AgentHandler {
	return &AgentHandler{s: sv}
}

func (h *AgentHandler) Create(gctx *gin.Context) {
}


func (h *AgentHandler) GetByID(gctx *gin.Context) error {
	return nil
}

func (h *AgentHandler) Fetch(gctx *gin.Context) error  {
	return nil
}

func (h *AgentHandler) FetchWithFilter(gctx *gin.Context) error  {
	return nil
}

func (h *AgentHandler) Update(gctx *gin.Context) error {
	return nil
}

func (h *AgentHandler) Delete(gctx *gin.Context) error {
	return nil
}
