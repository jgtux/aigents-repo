package handlers

import (
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"
	m "aigents-base/internal/auth-land/auth-signature/middleware"
	c_at "aigents-base/internal/common/atoms"
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
	authUUID, ok := m.GetAuthUUID(gctx)
	if !ok {
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusUnauthorized,
			"(H) Invalid context values.",
			"Invalid auth_uuid in context!")
		c_at.FeedErrLogToFile(err)
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Description string `json:"description"`
		ImageURL string `json:"image_url"`
		CategoryID uint64 `json:"category_id" binding:"required"`
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
		Name: req.Name,
		Description: req.Description,
		ImageURL: req.ImageURL,
		AuthUUID: authUUID,
	}
	agent.AgentConfig.Category.CategoryID = req.CategoryID
	// Temporaly
	agent.AgentConfig.CategoryPresetEnabled = false

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

	agent := &d.Agent{
		AgentUUID: req.AgentUUID.String(),
	}
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

