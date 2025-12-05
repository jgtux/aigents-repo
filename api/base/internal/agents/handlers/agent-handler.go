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

func NewAgentHandler(sv agitf.AgentServiceITF) *AgentHandler {
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

func (h *AgentHandler) GetByID(gctx *gin.Context) {
	param := gctx.Param("agent_uuid")
	agentUUID, err := uuid.Parse(param)
	if err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid URL parameter.",
			"Invalid agent_uuid param")
		c_at.FeedErrLogToFile(err)
		return
	}

	agent := &d.Agent{
		AgentUUID: agentUUID.String(),
	}

	err = h.s.GetByID(gctx, agent)
	if err != nil {
		c_at.FeedErrLogToFile(err)
	}

	data := d.Agent{
		AgentUUID: agent.AgentUUID,
		Name: agent.Name,
		Description: agent.Description,
		ImageURL: agent.ImageURL,
	}

	data.AgentConfig.Category.CategoryID = agent.AgentConfig.Category.CategoryID
	data.AgentConfig.Category.CategoryName = agent.AgentConfig.Category.CategoryName

	c_at.RespAtom[d.Agent](
		gctx,
		http.StatusOK,
		"(*) Data retrivied.",
		data)
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

	data, err := h.s.Fetch(gctx, req.PageSize, req.Page)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		return
	}

	c_at.RespAtom[[]d.Agent](gctx,
		http.StatusOK,
		"(*) Data retrivied",
		data)
}

func (h *AgentHandler) FetchByLoggedAuth(gctx *gin.Context) {
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

	data, err := h.s.FetchAgentsByLoggedAuth(gctx, authUUID, req.PageSize, req.Page)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		return
	}

	c_at.RespAtom[[]d.Agent](gctx,
		http.StatusOK,
		"(*) Data retrivied",
		data)
}

func (h *AgentHandler) FetchCategories(gctx *gin.Context) {
	data, err := h.s.FetchCategories(gctx)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		return
	}

	c_at.RespAtom[[]d.AgentCategory](gctx,
		http.StatusOK,
		"(*) Data retrivied",
		data)
}
