
package handlers

import (
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"
	c_at "aigents-base/internal/common/atoms"
	m "aigents-base/internal/auth-land/auth-signature/middleware"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	s agitf.AgentServiceITF
}

func NewAuthHandler(sv agitf.AgentServiceITF) *AgentHandler {
	return &AgentHandler{s: sv}
}

func (h *AgentHandler) Create(gctx *gin.Context) {
	var req struct {
		//
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid body request or values.",
			"Invalid body request")
		c_at.FeedErrLogToFile(err)
		return
	}

	agent := &d.Agent{
		
	}

	err := h.s.Create(gctx, agent)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		return
	}

	c_at.RespAtom[*struct{}](gctx,
		http.StatusCreated,
		"(*) Agent created",
		nil)

}


func (h *AgentHandler) GetByID(gctx *gin.Context)  {
	var req struct {
		AgentUUID uuid.UUID `json:"agent_uuid" binding:"required"`
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid body request or values.",
			"Invalid body request")
		c_at.FeedErrLogToFile(err)
		return
	}

	agent := &d.Agent{}
	err := h.s.GetByID(gctx, agent)
	if err != nil {
		c_at.FeedErrLogToFile(err)
	}

	c_at.RespAtom[d.Agent](
		gctx,
		http.StatusOK,
		"(*) Data retrivied.",
		*agent)
}

func (h *AgentHandler) Fetch(gctx *gin.Context) {
	var req struct {
		Page uint64 `json:"page"`
		PageSize uint64 `json:"page_size"`
	}

	err := gctx.ShouldBindJSON(&req)
	if  err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid body request or values.",
			"Invalid body request")
		c_at.FeedErrLogToFile(err)
		return
	}

	data, err := h.s.Fetch(gctx, req.Page, req.PageSize)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		return
	}

	c_at.RespAtom[[]d.Agent](gctx,
		http.StatusOK,
		"(*) Data retrivied",
		data)
}

