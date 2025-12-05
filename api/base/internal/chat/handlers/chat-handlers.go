package handlers

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	m "aigents-base/internal/auth-land/auth-signature/middleware"
	c_at "aigents-base/internal/common/atoms"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	s chitf.ChatServiceITF
}

func NewChatHandler(sv chitf.ChatServiceITF) *ChatHandler {
	return &ChatHandler{s: sv}
}

func (h *ChatHandler) Create(gctx *gin.Context) {
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
		AgentUUID      string `json:"agent_uuid" binding:"required"`
		MessageContent string `json:"message_content" binding:"required"`
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

	// Set up SSE headers for streaming
	gctx.Header("Content-Type", "text/event-stream")
	gctx.Header("Cache-Control", "no-cache")
	gctx.Header("Connection", "keep-alive")
	gctx.Header("Transfer-Encoding", "chunked")
	gctx.Header("X-Accel-Buffering", "no")

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(H) Streaming not supported.",
			"Streaming not supported")
		c_at.FeedErrLogToFile(err)
		return
	}

	// Create chat with initial message
	chatUUID := uuid.New().String()
	
	chat := &d.Chat{
		ChatUUID:  chatUUID,
		AuthUUID:  authUUID,
		AgentUUID: req.AgentUUID,
		History: []d.Message{
			{
				MessageUUID:  uuid.New().String(),
				SenderUUID:   authUUID,
				SenderType:   "AUTH",
				ReceiverUUID: req.AgentUUID,
				ReceiverType: "AGENT",
				MessageContent: d.MessageContent{
					MessageContentUUID: uuid.New().String(),
					Content:            req.MessageContent,
				},
				CreatedAt: time.Now(),
			},
		},
	}

	// Send initial test event
	gctx.SSEvent("test", "connection established")
	flusher.Flush()

	// Stream callback function
	streamCallback := func(chunk string) {
		gctx.SSEvent("message", chunk)
		flusher.Flush()
	}

	// Call service with streaming
	err := h.s.InitChat(gctx, chat, streamCallback)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		gctx.SSEvent("error", "(SSE) Could not initialize chat.")
		flusher.Flush()
		return
	}

	// Send final message with complete chat
	gctx.SSEvent("done", chat)
	flusher.Flush()
}

func (h *ChatHandler) SendMessage(gctx *gin.Context) {
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
		ChatUUID       string `json:"chat_uuid" binding:"required"`
		MessageContent string `json:"message_content" binding:"required"`
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

	// Build user message
	userMessage := &d.Message{
		MessageUUID: uuid.New().String(),
		ChatUUID:    req.ChatUUID,
		SenderUUID:  authUUID,
		SenderType:  "AUTH",
		MessageContent: d.MessageContent{
			MessageContentUUID: uuid.New().String(),
			Content:            req.MessageContent,
		},
		CreatedAt: time.Now(),
	}

	// Set up SSE headers for streaming
	gctx.Header("Content-Type", "text/event-stream")
	gctx.Header("Cache-Control", "no-cache")
	gctx.Header("Connection", "keep-alive")
	gctx.Header("Transfer-Encoding", "chunked")
	gctx.Header("X-Accel-Buffering", "no")

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(H) Streaming not supported.",
			"Streaming not supported")
		c_at.FeedErrLogToFile(err)
		return
	}

	// Send initial test event
	gctx.SSEvent("test", "connection established")
	flusher.Flush()

	// Stream callback function
	streamCallback := func(chunk string) {
		gctx.SSEvent("message", chunk)
		flusher.Flush()
	}

	// Call service with streaming
	err := h.s.SendMessage(gctx, userMessage, authUUID, streamCallback)
	if err != nil {
		c_at.FeedErrLogToFile(err)
		gctx.SSEvent("error", "(SSE) Could not send message.")
		flusher.Flush()
		return
	}

	// Send final message with complete response
	gctx.SSEvent("done", userMessage)
	flusher.Flush()
}
