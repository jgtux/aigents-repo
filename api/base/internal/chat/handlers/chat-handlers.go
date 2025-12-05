package handlers

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	m "aigents-base/internal/auth-land/auth-signature/middleware"
	c_at "aigents-base/internal/common/atoms"
	"fmt"
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
	fmt.Println("[DEBUG] Chat Create endpoint called")
	
	authUUID, ok := m.GetAuthUUID(gctx)
	if !ok {
		fmt.Println("[DEBUG] Failed to get authUUID from context")
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusUnauthorized,
			"(H) Invalid context values.",
			"Invalid auth_uuid in context!")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Printf("[DEBUG] AuthUUID: %s\n", authUUID)

	var req struct {
		AgentUUID      string `json:"agent_uuid" binding:"required"`
		MessageContent string `json:"message_content" binding:"required"`
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[DEBUG] Failed to bind JSON: %v\n", err)
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid body request or values.",
			"Invalid body request")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Printf("[DEBUG] Request - AgentUUID: %s, MessageContent: %s\n", req.AgentUUID, req.MessageContent)

	// Set up SSE headers for streaming
	fmt.Println("[DEBUG] Setting up SSE headers")
	gctx.Header("Content-Type", "text/event-stream")
	gctx.Header("Cache-Control", "no-cache")
	gctx.Header("Connection", "keep-alive")
	gctx.Header("Transfer-Encoding", "chunked")
	gctx.Header("X-Accel-Buffering", "no") // Disable nginx buffering if applicable

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		fmt.Println("[DEBUG] Streaming not supported - flusher unavailable")
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(H) Streaming not supported.",
			"Streaming not supported")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Println("[DEBUG] Flusher obtained successfully")

	// Create chat with initial message
	chatUUID := uuid.New().String()
	fmt.Printf("[DEBUG] Generated ChatUUID: %s\n", chatUUID)
	
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
	fmt.Println("[DEBUG] Sending initial test event")
	gctx.SSEvent("test", "connection established")
	flusher.Flush()

	// Stream callback function
	streamCallback := func(chunk string) {
		fmt.Printf("[DEBUG] Streaming chunk: %s\n", chunk)
		gctx.SSEvent("message", chunk)
		flusher.Flush()
	}

	// Call service with streaming
	fmt.Println("[DEBUG] Calling InitChat service")
	err := h.s.InitChat(gctx, chat, streamCallback)
	if err != nil {
		fmt.Printf("[DEBUG] InitChat error: %v\n", err)
		gctx.SSEvent("error", err.Error())
		flusher.Flush()
		return
	}
	fmt.Println("[DEBUG] InitChat completed successfully")

	// Send final message with complete chat
	fmt.Printf("[DEBUG] Sending done event with ChatUUID: %s\n", chat.ChatUUID)
	gctx.SSEvent("done", chat)
	flusher.Flush()
	fmt.Println("[DEBUG] Chat Create handler completed")
}

func (h *ChatHandler) SendMessage(gctx *gin.Context) {
	fmt.Println("[DEBUG] SendMessage endpoint called")
	
	authUUID, ok := m.GetAuthUUID(gctx)
	if !ok {
		fmt.Println("[DEBUG] Failed to get authUUID from context")
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusUnauthorized,
			"(H) Invalid context values.",
			"Invalid auth_uuid in context!")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Printf("[DEBUG] AuthUUID: %s\n", authUUID)

	var req struct {
		ChatUUID       string `json:"chat_uuid" binding:"required"`
		MessageContent string `json:"message_content" binding:"required"`
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[DEBUG] Failed to bind JSON: %v\n", err)
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusBadRequest,
			"(H) Invalid body request or values.",
			"Invalid body request")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Printf("[DEBUG] Request - ChatUUID: %s, MessageContent: %s\n", req.ChatUUID, req.MessageContent)

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
	fmt.Printf("[DEBUG] Created user message with UUID: %s\n", userMessage.MessageUUID)

	// Set up SSE headers for streaming
	fmt.Println("[DEBUG] Setting up SSE headers")
	gctx.Header("Content-Type", "text/event-stream")
	gctx.Header("Cache-Control", "no-cache")
	gctx.Header("Connection", "keep-alive")
	gctx.Header("Transfer-Encoding", "chunked")
	gctx.Header("X-Accel-Buffering", "no") // Disable nginx buffering if applicable

	flusher, ok := gctx.Writer.(http.Flusher)
	if !ok {
		fmt.Println("[DEBUG] Streaming not supported - flusher unavailable")
		err := c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(H) Streaming not supported.",
			"Streaming not supported")
		c_at.FeedErrLogToFile(err)
		return
	}
	fmt.Println("[DEBUG] Flusher obtained successfully")

	// Send initial test event
	fmt.Println("[DEBUG] Sending initial test event")
	gctx.SSEvent("test", "connection established")
	flusher.Flush()

	// Stream callback function
	streamCallback := func(chunk string) {
		fmt.Printf("[DEBUG] Streaming chunk: %s\n", chunk)
		gctx.SSEvent("message", chunk)
		flusher.Flush()
	}

	// Call service with streaming
	fmt.Println("[DEBUG] Calling SendMessage service")
	err := h.s.SendMessage(gctx, userMessage, authUUID, streamCallback)
	if err != nil {
		fmt.Printf("[DEBUG] SendMessage error: %v\n", err)
		gctx.SSEvent("error", err.Error())
		flusher.Flush()
		return
	}
	fmt.Println("[DEBUG] SendMessage completed successfully")

	// Send final message with complete response
	fmt.Printf("[DEBUG] Sending done event with MessageUUID: %s\n", userMessage.MessageUUID)
	gctx.SSEvent("done", userMessage)
	flusher.Flush()
	fmt.Println("[DEBUG] SendMessage handler completed")
}
