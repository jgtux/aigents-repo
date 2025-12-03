package services

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	agitf "aigents-base/internal/agents/interfaces"
	"fmt"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)


type PythonLLMRequest struct {
	ChatUUID          string      `json:"chat_uuid"`
	Content           string      `json:"content"`
	SenderUUID        string      `json:"sender_uuid"`
	SenderType        string      `json:"sender_type"`
	ReceiverUUID      string      `json:"receiver_uuid"`
	ReceiverType      string      `json:"receiver_type"`
	AgentUUID         string      `json:"agent_uuid"`
	AgentName         string      `json:"agent_name"`
	AgentDescription  string      `json:"agent_description"`
	CategoryID        uint64      `json:"category_id"`
	SystemPrompt      string      `json:"system_prompt"`
	ChatHistory       []d.Message `json:"chat_history,omitempty"`
	SyncMode          string      `json:"sync_mode"` // "auto", "incremental", "full"
}

type PythonLLMResponse struct {
	ChatUUID           string `json:"chat_uuid"`
	AgentUUID          string `json:"agent_uuid"`
	Content            string `json:"content"`
	Partial            bool   `json:"partial"`
	MessageUUID        string `json:"message_uuid,omitempty"`
	MessageContentUUID string `json:"message_content_uuid,omitempty"`
	Error              string `json:"error,omitempty"`
}


// Connection pool manages multiple WebSocket connections
type ConnectionPool struct {
	wsURL     string
	pool      chan *ws.Conn
	mu        sync.Mutex
	maxConns  int
	connCount int
}

func NewConnectionPool(wsURL string, size int) *ConnectionPool {
	return &ConnectionPool{
		wsURL:    wsURL,
		pool:     make(chan *ws.Conn, size),
		maxConns: size,
	}
}

// Get a connection from pool or create new one
func (p *ConnectionPool) Get() (*ws.Conn, error) {
	select {
	case conn := <-p.pool:
		// Check if connection is still alive
		if err := conn.WriteControl(ws.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
			conn.Close()
			return p.createConnection()
		}
		return conn, nil
	default:
		// Pool empty, create new if under limit
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.connCount < p.maxConns {
			return p.createConnection()
		}

		// Wait for available connection
		return <-p.pool, nil
	}
}

// Return connection to pool
func (p *ConnectionPool) Put(conn *ws.Conn) {
	if conn == nil {
		return
	}

	select {
	case p.pool <- conn:
		// Successfully returned to pool
	default:
		// Pool full, close connection
		conn.Close()
		p.mu.Lock()
		p.connCount--
		p.mu.Unlock()
	}
}

func (p *ConnectionPool) createConnection() (*ws.Conn, error) {
	conn, _, err := ws.DefaultDialer.Dial(p.wsURL, nil)
	if err != nil {
		return nil, err
	}
	p.connCount++
	return conn, nil
}

func (p *ConnectionPool) Close() {
	close(p.pool)
	for conn := range p.pool {
		conn.Close()
	}
}

// ChatService with connection pool
type ChatService struct {
	r             chitf.ChatRepositoryITF
	agr           agitf.AgentRepositoryITF
	lastMsgsLimit uint64
	connPool      *ConnectionPool
}

func NewChatService(repo chitf.ChatRepositoryITF, agrepo agitf.AgentRepositoryITF, wsURL string, lastMsgsLimit uint64, poolSize int) chitf.ChatServiceITF {
	return &ChatService{
		r:             repo,
		lastMsgsLimit: lastMsgsLimit,
		connPool:      NewConnectionPool(wsURL, poolSize),
	}
}

func (s *ChatService) SendMessage(gctx *gin.Context, data *d.Message, authUUID string, streamCallback func(chunk string)) error {
	// Get connection from pool
	conn, err := s.connPool.Get()
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	// Connection is exclusive to this goroutine now
	// Return it when done (not in defer to handle errors properly)
	returnConn := true
	defer func() {
		if returnConn {
			s.connPool.Put(conn)
		} else {
			conn.Close()
		}
	}()

	chat := &d.Chat{ ChatUUID: data.ChatUUID }
	err = s.r.GetByID(gctx, chat)
	if err != nil {
		return err
	}

	// 2. Get agent configuration
	agent, err := s.agr.GetAgentByUUID(gctx, chat.AgentUUID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	// 3. Get chat history
	err = s.r.GetChatHistory(gctx, chat, s.lastMsgsLimit)
	if err != nil {
		return fmt.Errorf("failed to get chat history: %w", err)
	}

	// 4. Determine sync mode
	syncMode := s.determineChatHistoryStrategy(chat, uint64(len(chat.History)))

	// 5. Save user message to DB first
	if err := s.r.AttachMessage(gctx, data); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// 6. Extract system prompt
	systemPrompt := "You are a helpful assistant."
	if agent.AgentConfig.AgentSystem.SystemPreset != nil {
		if prompt, ok := agent.AgentConfig.AgentSystem.SystemPreset["system_prompt"].(string); ok {
			systemPrompt = prompt
		}
	}

	// 7. Build request for Python
	request := PythonLLMRequest{
		ChatUUID:         data.ChatUUID,
		Content:          data.MessageContent.Content,
		SenderUUID:       authUUID,
		SenderType:       "AUTH",
		ReceiverUUID:     agent.AgentUUID,
		ReceiverType:     "AGENT",
		AgentUUID:        agent.AgentUUID,
		AgentName:        agent.Name,
		AgentDescription: agent.Description,
		CategoryID:       1,
		SystemPrompt:     systemPrompt,
		ChatHistory:      chat.History,
		SyncMode:         syncMode,
	}

	// 8. Send request to Python
	if err := conn.WriteJSON(request); err != nil {
		returnConn = false // Connection broken, don't return to pool
		return fmt.Errorf("failed to send request: %w", err)
	}

	// 9. Stream response chunks
	var agentMessageUUID string
	var messageContent d.MessageContent

	for {
		var response PythonLLMResponse
		err = conn.ReadJSON(&response)
		if err != nil {
			returnConn = false // Connection broken
			return fmt.Errorf("failed to read response: %w", err)
		}

		if response.Error != "" {
			return fmt.Errorf("python service error: %s", response.Error)
		}

		if response.Partial {
			// Stream chunk to client
			if streamCallback != nil {
				streamCallback(response.Content)
			}
		} else {
			// Final response
			messageContent.Content = response.Content
			messageContent.MessageContentUUID = response.MessageContentUUID
			agentMessageUUID = response.MessageUUID
			break
		}
	}

	// 10. Save agent response to DB
	agentMsg := &d.Message{
		MessageUUID:    agentMessageUUID,
		SenderUUID:     agent.AgentUUID,
		SenderType:     "AGENT",
		ReceiverUUID:   authUUID,
		ReceiverType:   "AUTH",
		ChatUUID:       chat.ChatUUID,
		MessageContent: messageContent,
		CreatedAt:      time.Now(),
	}

	if err := s.r.AttachMessage(gctx, agentMsg); err != nil {
		return fmt.Errorf("failed to save agent message: %w", err)
	}

	*data = *agentMsg
	return nil
}

func (s *ChatService) determineChatHistoryStrategy(data *d.Chat, msgsLen uint64) string {
	if time.Since(data.UpdatedAt) < 5*time.Minute {
		return "auto"
	}
	if time.Since(data.UpdatedAt) > 1*time.Hour || msgsLen > 15 {
		return "full"
	}
	return "auto"
}


func (s *ChatService) Create(gctx *gin.Context, data *d.Chat) error { return nil }
func (s *ChatService) GetByID(gctx *gin.Context, data *d.Chat) error { return s.r.GetByID(gctx, data) }
func (s *ChatService) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Chat, error) { return s.r.Fetch(gctx, limit, offset) }
func (s *ChatService) Update(gctx *gin.Context, data *d.Chat) error { return s.r.Update(gctx, data) }
func (s *ChatService) Delete(gctx *gin.Context, data *d.Chat) error { return s.r.Delete(gctx, data) }
