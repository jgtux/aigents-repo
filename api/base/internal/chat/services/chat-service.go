package services

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	agitf "aigents-base/internal/agents/interfaces"
	c_at "aigents-base/internal/common/atoms"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type PythonLLMRequest struct {
	Command          string      `json:"command,omitempty"`
	ChatUUID         string      `json:"chat_uuid"`
	Content          string      `json:"content"`
	SenderUUID       string      `json:"sender_uuid"`
	SenderType       string      `json:"sender_type"`
	ReceiverUUID     string      `json:"receiver_uuid"`
	ReceiverType     string      `json:"receiver_type"`
	AgentUUID        string      `json:"agent_uuid"`
	AgentName        string      `json:"agent_name"`
	AgentDescription string      `json:"agent_description"`
	CategoryID       uint64      `json:"category_id"`
	SystemPrompt     string      `json:"system_prompt"`
	ChatHistory      []d.Message `json:"chat_history,omitempty"`
	SyncMode         string      `json:"sync_mode"`
	AuthUUID         string      `json:"auth_uuid,omitempty"`
}

type PythonLLMResponse struct {
	Type               string `json:"type,omitempty"`
	ConnectionID       string `json:"connection_id,omitempty"`
	ChatUUID           string `json:"chat_uuid"`
	AgentUUID          string `json:"agent_uuid"`
	Content            string `json:"content"`
	Partial            bool   `json:"partial"`
	MessageUUID        string `json:"message_uuid,omitempty"`
	MessageContentUUID string `json:"message_content_uuid,omitempty"`
	Error              string `json:"error,omitempty"`
}

type PooledConnection struct {
	Conn         *ws.Conn
	ConnectionID string
	CreatedAt    time.Time
	LastUsed     time.Time
	UseCount     int
	Identified   bool
}

type ConnectionPool struct {
	wsURL           string
	pool            chan *PooledConnection
	mu              sync.RWMutex
	maxConns        int
	activeConns     map[string]*PooledConnection
	connTimeout     time.Duration
	maxConnAge      time.Duration
	maxUseCount     int
	identifyTimeout time.Duration
	closed          bool
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewConnectionPool(wsURL string, size int) *ConnectionPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &ConnectionPool{
		wsURL:           wsURL,
		pool:            make(chan *PooledConnection, size),
		maxConns:        size,
		activeConns:     make(map[string]*PooledConnection),
		connTimeout:     30 * time.Second,
		maxConnAge:      10 * time.Minute,
		maxUseCount:     100,
		identifyTimeout: 5 * time.Second,
		ctx:             ctx,
		cancel:          cancel,
	}

	go pool.cleanupStaleConnections()

	return pool
}

func (p *ConnectionPool) Get(authUUID string) (*PooledConnection, error) {
	if p.closed {
		return nil, fmt.Errorf("connection pool is closed")
	}

	select {
	case pooledConn := <-p.pool:
		if p.isHealthy(pooledConn) {
			pooledConn.LastUsed = time.Now()
			pooledConn.UseCount++

			p.mu.Lock()
			p.activeConns[pooledConn.ConnectionID] = pooledConn
			p.mu.Unlock()

			return pooledConn, nil
		}

		pooledConn.Conn.Close()

	default:
	}

	p.mu.Lock()
	currentCount := len(p.activeConns)
	if currentCount >= p.maxConns {
		p.mu.Unlock()
		select {
		case pooledConn := <-p.pool:
			if p.isHealthy(pooledConn) {
				pooledConn.LastUsed = time.Now()
				pooledConn.UseCount++
				p.mu.Lock()
				p.activeConns[pooledConn.ConnectionID] = pooledConn
				p.mu.Unlock()
				return pooledConn, nil
			}
			pooledConn.Conn.Close()
			return nil, fmt.Errorf("all connections are unhealthy")
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("timeout waiting for available connection")
		}
	}
	p.mu.Unlock()

	return p.createConnection(authUUID)
}

func (p *ConnectionPool) createConnection(authUUID string) (*PooledConnection, error) {
	dialer := ws.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(p.wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	connectionID := fmt.Sprintf("go-%d", time.Now().UnixNano())

	pooledConn := &PooledConnection{
		Conn:         conn,
		ConnectionID: connectionID,
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		UseCount:     0,
		Identified:   false,
	}

	if err := p.identifyConnection(pooledConn, authUUID); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to identify connection: %w", err)
	}

	p.mu.Lock()
	p.activeConns[connectionID] = pooledConn
	p.mu.Unlock()

	return pooledConn, nil
}

func (p *ConnectionPool) identifyConnection(pooledConn *PooledConnection, authUUID string) error {
	identifyMsg := PythonLLMRequest{
		Command:  "identify",
		AuthUUID: authUUID,
	}

	pooledConn.Conn.SetWriteDeadline(time.Now().Add(p.identifyTimeout))
	if err := pooledConn.Conn.WriteJSON(identifyMsg); err != nil {
		return fmt.Errorf("failed to send identify: %w", err)
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})

	pooledConn.Conn.SetReadDeadline(time.Now().Add(p.identifyTimeout))
	var response PythonLLMResponse
	if err := pooledConn.Conn.ReadJSON(&response); err != nil {
		return fmt.Errorf("failed to read identify response: %w", err)
	}
	pooledConn.Conn.SetReadDeadline(time.Time{})

	if response.Type != "identified" {
		return fmt.Errorf("unexpected response type: %s", response.Type)
	}

	pooledConn.Identified = true

	return nil
}

func (p *ConnectionPool) isHealthy(pooledConn *PooledConnection) bool {
	if time.Since(pooledConn.CreatedAt) > p.maxConnAge {
		return false
	}

	if pooledConn.UseCount >= p.maxUseCount {
		return false
	}

	deadline := time.Now().Add(2 * time.Second)
	if err := pooledConn.Conn.WriteControl(ws.PingMessage, []byte{}, deadline); err != nil {
		return false
	}

	return true
}

func (p *ConnectionPool) Put(pooledConn *PooledConnection) {
	if pooledConn == nil || p.closed {
		return
	}

	p.mu.Lock()
	delete(p.activeConns, pooledConn.ConnectionID)
	p.mu.Unlock()

	if !p.isHealthy(pooledConn) {
		pooledConn.Conn.Close()
		return
	}

	select {
	case p.pool <- pooledConn:
	default:
		pooledConn.Conn.Close()
	}
}

func (p *ConnectionPool) Discard(pooledConn *PooledConnection) {
	if pooledConn == nil {
		return
	}

	p.mu.Lock()
	delete(p.activeConns, pooledConn.ConnectionID)
	p.mu.Unlock()

	pooledConn.Conn.Close()
}

func (p *ConnectionPool) cleanupStaleConnections() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.mu.Lock()
			now := time.Now()
			stale := make([]string, 0)

			for id, conn := range p.activeConns {
				if now.Sub(conn.LastUsed) > 5*time.Minute {
					stale = append(stale, id)
				}
			}

			for _, id := range stale {
				conn := p.activeConns[id]
				delete(p.activeConns, id)
				conn.Conn.Close()
			}
			p.mu.Unlock()

			poolSize := len(p.pool)
			for i := 0; i < poolSize; i++ {
				select {
				case conn := <-p.pool:
					if !p.isHealthy(conn) {
						conn.Conn.Close()
					} else {
						select {
						case p.pool <- conn:
						default:
							conn.Conn.Close()
						}
					}
				default:
					break
				}
			}
		}
	}
}

func (p *ConnectionPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	p.cancel()

	p.mu.Lock()
	for id, conn := range p.activeConns {
		conn.Conn.Close()
		delete(p.activeConns, id)
	}
	p.mu.Unlock()

	close(p.pool)
	for conn := range p.pool {
		conn.Conn.Close()
	}
}

func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"max_connections":    p.maxConns,
		"active_connections": len(p.activeConns),
		"pooled_connections": len(p.pool),
		"pool_capacity":      cap(p.pool),
	}
}

type ChatService struct {
	r             chitf.ChatRepositoryITF
	agr           agitf.AgentRepositoryITF
	lastMsgsLimit uint64
	connPool      *ConnectionPool
}

func NewChatService(repo chitf.ChatRepositoryITF, agrepo agitf.AgentRepositoryITF, wsURL string, lastMsgsLimit uint64, poolSize int) chitf.ChatServiceITF {
	return &ChatService{
		r:             repo,
		agr:           agrepo,
		lastMsgsLimit: lastMsgsLimit,
		connPool:      NewConnectionPool(wsURL, poolSize),
	}
}

func (s *ChatService) SendMessage(gctx *gin.Context, data *d.Message, authUUID string, streamCallback func(chunk string)) error {
	pooledConn, err := s.connPool.Get(authUUID)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not connect to AI service. Failed to get connection: %s", err.Error()))
		return err
	}

	shouldReturn := true
	defer func() {
		if shouldReturn {
			s.connPool.Put(pooledConn)
		} else {
			s.connPool.Discard(pooledConn)
		}
	}()

	chat := &d.Chat{ChatUUID: data.ChatUUID}
	if err := s.r.GetByID(gctx, chat); err != nil {
		return err
	}

	agent, err := s.agr.GetAgentByUUID(gctx, chat.AgentUUID)
	if err != nil {
		return err
	}

	data.ReceiverUUID = agent.AgentUUID
	data.ReceiverType = "AGENT"

	if err := s.r.AttachMessage(gctx, data); err != nil {
		return err
	}

	chat.History, err = s.r.GetChatHistory(gctx, chat.ChatUUID, s.lastMsgsLimit+1)
	if err != nil {
		return err
	}

	historyForPython := make([]d.Message, 0, len(chat.History))
	for _, msg := range chat.History {
		if msg.MessageUUID != data.MessageUUID {
			historyForPython = append(historyForPython, msg)
		}
	}

	syncMode := s.determineChatHistoryStrategy(chat, uint64(len(historyForPython)))

	systemPrompt := "You are a helpful assistant."
	if agent.AgentConfig.AgentSystem.SystemPreset != nil {
		if prompt, ok := agent.AgentConfig.AgentSystem.SystemPreset["system_prompt"].(string); ok {
			systemPrompt = prompt
		}
	}

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
		ChatHistory:      historyForPython,
		SyncMode:         syncMode,
	}

	pooledConn.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	if err := pooledConn.Conn.WriteJSON(request); err != nil {
		shouldReturn = false
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("AI service is unavailable. Failed to send request to Python service: %s", err.Error()))
		return err
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})

	var agentMessageUUID string
	var messageContent d.MessageContent

	for {
		pooledConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var response PythonLLMResponse
		err = pooledConn.Conn.ReadJSON(&response)
		if err != nil {
			shouldReturn = false
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("Failed to receive AI response. Failed to read response from Python service: %s", err.Error()))
			return err
		}

		if response.Error != "" {
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("AI service encountered an error. Python service error: %s", response.Error))
			return err
		}

		if response.Partial {
			if streamCallback != nil {
				streamCallback(response.Content)
			}
		} else {
			messageContent.Content = response.Content
			messageContent.MessageContentUUID = response.MessageContentUUID
			agentMessageUUID = response.MessageUUID
			break
		}
	}

	pooledConn.Conn.SetReadDeadline(time.Time{})

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
		return err
	}

	*data = *agentMsg
	return nil
}

func (s *ChatService) InitChat(gctx *gin.Context, data *d.Chat, streamCallback func(chunk string)) error {
	if len(data.History) == 0 {
		err := c_at.BuildErrLogAtom(
			gctx,
			"At least one message is required.")
		return err
	}

	pooledConn, err := s.connPool.Get(data.AuthUUID)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not connect to AI service. Failed to get connection: %s", err.Error()))
		return err
	}

	shouldReturn := true
	defer func() {
		if shouldReturn {
			s.connPool.Put(pooledConn)
		} else {
			s.connPool.Discard(pooledConn)
		}
	}()

	now := time.Now()
	if data.CreatedAt.IsZero() {
		data.CreatedAt = now
	}
	if data.UpdatedAt.IsZero() {
		data.UpdatedAt = now
	}

	if err := s.r.Create(gctx, data); err != nil {
		return err
	}

	agent, err := s.agr.GetAgentByUUID(gctx, data.AgentUUID)
	if err != nil {
		return err
	}

	userMessage := data.History[0]
	userMessage.ChatUUID = data.ChatUUID

	if userMessage.CreatedAt.IsZero() {
		userMessage.CreatedAt = now
	}

	if err := s.r.AttachMessage(gctx, &userMessage); err != nil {
		return err
	}

	systemPrompt := "You are a helpful assistant."
	if agent.AgentConfig.AgentSystem.SystemPreset != nil {
		if prompt, ok := agent.AgentConfig.AgentSystem.SystemPreset["system_prompt"].(string); ok {
			systemPrompt = prompt
		}
	}

	request := PythonLLMRequest{
		ChatUUID:         data.ChatUUID,
		Content:          userMessage.MessageContent.Content,
		SenderUUID:       userMessage.SenderUUID,
		SenderType:       userMessage.SenderType,
		ReceiverUUID:     agent.AgentUUID,
		ReceiverType:     "AGENT",
		AgentUUID:        agent.AgentUUID,
		AgentName:        agent.Name,
		AgentDescription: agent.Description,
		CategoryID:       1,
		SystemPrompt:     systemPrompt,
		ChatHistory:      []d.Message{},
		SyncMode:         "auto",
	}

	pooledConn.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	if err := pooledConn.Conn.WriteJSON(request); err != nil {
		shouldReturn = false
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("AI service is unavailable. Failed to send request to Python service: %s", err.Error()))
		return err
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})

	var agentMessageUUID string
	var messageContent d.MessageContent

	for {
		pooledConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var response PythonLLMResponse
		err = pooledConn.Conn.ReadJSON(&response)
		if err != nil {
			shouldReturn = false
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("Failed to receive AI response. Failed to read response from Python service: %s", err.Error()))
			return err
		}

		if response.Error != "" {
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("AI service encountered an error. Python service error: %s", response.Error))
			return err
		}

		if response.Partial {
			if streamCallback != nil {
				streamCallback(response.Content)
			}
		} else {
			messageContent.Content = response.Content
			messageContent.MessageContentUUID = response.MessageContentUUID
			agentMessageUUID = response.MessageUUID
			break
		}
	}

	pooledConn.Conn.SetReadDeadline(time.Time{})

	agentMsg := &d.Message{
		MessageUUID:    agentMessageUUID,
		SenderUUID:     agent.AgentUUID,
		SenderType:     "AGENT",
		ReceiverUUID:   userMessage.SenderUUID,
		ReceiverType:   userMessage.SenderType,
		ChatUUID:       data.ChatUUID,
		MessageContent: messageContent,
		CreatedAt:      time.Now(),
	}

	if err := s.r.AttachMessage(gctx, agentMsg); err != nil {
		return err
	}

	data.History = []d.Message{userMessage, *agentMsg}

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

func (s *ChatService) Cleanup() {
	s.connPool.Close()
}
