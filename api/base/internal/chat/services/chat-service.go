 package services

import (
	"fmt"
	"encoding/json"
	"time"
	agitf "aigents-base/internal/agents/interfaces"
	chitf "aigents-base/internal/chat/interfaces"
	ws "github.com/gorilla/websocket"
	"net/http"
	d "aigents-base/internal/chat/domain"
	"github.com/gin-gonic/gin"
)


type PythonLLMRequest struct {
	ChatUUID          string              `json:"chat_uuid"`
	Content           string                 `json:"content"`
	SenderUUID        string              `json:"sender_uuid"`
	SenderType        string             `json:"sender_type"`
	ReceiverUUID      string              `json:"receiver_uuid"`
	ReceiverType      string             `json:"receiver_type"`
	AgentUUID         string              `json:"agent_uuid"`
	AgentName         string                 `json:"agent_name"`
	AgentDescription  string                 `json:"agent_description"`
	CategoryID        uint64                    `json:"category_id"`
	SystemPrompt      string                 `json:"system_prompt"`
	ChatHistory       []d.Message              `json:"chat_history,omitempty"`
	SyncMode          string                 `json:"sync_mode"` // "auto", "incremental", "full"
}

type PythonLLMResponse struct {
	ChatUUID           string `json:"chat_uuid"`
	AgentUUID          string `json:"agent_uuid"`
	Content            string    `json:"content"`
	Partial            bool      `json:"partial"`
	MessageUUID        string `json:"message_uuid,omitempty"`
	MessageContentUUID string `json:"message_content_uuid,omitempty"`
	Error              string    `json:"error,omitempty"`
}

type ChatService struct {
	ws_str string
	r chitf.ChatRepositoryITF
	ar agitf.AgentRepositoryITF
}

func NewChatService(repo chitf.ChatRepositoryITF, ws_str string) chitf.ChatServiceITF {
	return &ChatService{r: repo, ws_str: ws_str}
}

func (s *ChatService) Create(gctx *gin.Context, data *d.Chat) error {
	return nil
}

func (s *ChatService) GetByID(gctx *gin.Context, data *d.Chat) error {
	return s.r.GetByID(gctx, data)
}

func (s *ChatService) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Chat, error) {
	return s.r.Fetch(gctx, limit, offset)
}

func (s *ChatService) Update(gctx *gin.Context, data *d.Chat) error {
	return s.r.Update(gctx, data)
}

func (s *ChatService) Delete(gctx *gin.Context, data *d.Chat) error {
	return s.r.Delete(gctx, data)
}


func (s *ChatService) SendMessage(gctx *gin.Context, data *d.Message, authUUID string, streamCallback func(chunk string)) error {
	// // 1. Get chat metadata
	// chat, err := s.chatService.GetChatMetadata(ctx, chatUUID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get chat: %w", err)
	// }

	// // 2. Get agent configuration
	// agent, agentSystem, err := s.chatService.GetAgentWithSystem(ctx, chat.AgentUUID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get agent: %w", err)
	// }

	// // 3. Determine if we need to send chat history
	// chatHistory, syncMode, err := s.determineChatHistoryStrategy(ctx, chatUUID, chat.UpdatedAt)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to determine history strategy: %w", err)
	// }

	// // 4. Save user message to DB first
	// userMsg := Message{
	// 	MessageUUID:        uuid.New(),
	// 	SenderUUID:         authUUID,
	// 	SenderType:         EntityTypeAuth,
	// 	ReceiverUUID:       agent.AgentUUID,
	// 	ReceiverType:       EntityTypeAgent,
	// 	ChatUUID:           chatUUID,
	// 	MessageContentUUID: uuid.New(),
	// 	Content:            content,
	// 	CreatedAt:          time.Now(),
	// }

	// if err := s.chatService.SaveMessage(ctx, userMsg); err != nil {
	// 	return nil, fmt.Errorf("failed to save user message: %w", err)
	// }

	// // 5. Connect to Python WebSocket
	// conn, _, err := websocket.DefaultDialer.Dial(s.pythonWSURL, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to Python service: %w", err)
	// }
	// defer conn.Close()

	// // 6. Extract system prompt from agent system
	// systemPrompt := "You are a helpful assistant."
	// if prompt, ok := agentSystem.SystemPreset["system_prompt"].(string); ok {
	// 	systemPrompt = prompt
	// }

	// // 7. Build request for Python
	// request := PythonLLMRequest{
	// 	ChatUUID:         chatUUID,
	// 	Content:          content,
	// 	SenderUUID:       authUUID,
	// 	SenderType:       EntityTypeAuth,
	// 	ReceiverUUID:     agent.AgentUUID,
	// 	ReceiverType:     EntityTypeAgent,
	// 	AgentUUID:        agent.AgentUUID,
	// 	AgentName:        agent.Name,
	// 	AgentDescription: agent.Description,
	// 	CategoryID:       1, // You might want to get this from agent_config
	// 	SystemPrompt:     systemPrompt,
	// 	ChatHistory:      chatHistory,
	// 	SyncMode:         syncMode,
	// }

	// // 8. Send request to Python
	// if err := conn.WriteJSON(request); err != nil {
	// 	return nil, fmt.Errorf("failed to send request to Python: %w", err)
	// }

	// // 9. Stream response chunks
	// var fullResponse string
	// var agentMessageUUID uuid.UUID
	// var agentMessageContentUUID uuid.UUID

	// for {
	// 	var response PythonLLMResponse
	// 	if err := conn.ReadJSON(&response); err != nil {
	// 		return nil, fmt.Errorf("failed to read Python response: %w", err)
	// 	}

	// 	if response.Error != "" {
	// 		return nil, fmt.Errorf("Python LLM error: %s", response.Error)
	// 	}

	// 	if response.Partial {
	// 		// Stream chunk to client
	// 		if streamCallback != nil {
	// 			streamCallback(response.Content)
	// 		}
	// 	} else {
	// 		// Final response
	// 		fullResponse = response.Content
	// 		agentMessageUUID = response.MessageUUID
	// 		agentMessageContentUUID = response.MessageContentUUID
	// 		break
	// 	}
	// }

	// // 10. Save agent response to DB
	// agentMsg := Message{
	// 	MessageUUID:        agentMessageUUID,
	// 	SenderUUID:         agent.AgentUUID,
	// 	SenderType:         EntityTypeAgent,
	// 	ReceiverUUID:       authUUID,
	// 	ReceiverType:       EntityTypeAuth,
	// 	ChatUUID:           chatUUID,
	// 	MessageContentUUID: agentMessageContentUUID,
	// 	Content:            fullResponse,
	// 	CreatedAt:          time.Now(),
	// }

	// if err := s.chatService.SaveMessage(ctx, agentMsg); err != nil {
	// 	return nil, fmt.Errorf("failed to save agent message: %w", err)
	// }

	// return &agentMsg, nil
}

func (s *ChatService) determineChatHistoryStrategy(data *d.Chat, msgsLen uint64) string {
	// Strategy 1: For very recent chats (< 5 min old), assume Python has cache
	if time.Since(data.UpdatedAt) < 5*time.Minute {
		return "auto"
	}

	// If chat is old (> 1 hour) or has many messages, force full reload
	if time.Since(data.UpdatedAt) > 1*time.Hour || msgsLen > 15 {
		return "full"
	}

	// Otherwise, let Python decide (auto mode)
	return "auto"

}
