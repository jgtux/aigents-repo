package services

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	agitf "aigents-base/internal/agents/interfaces"
	"context"
	"fmt"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type PythonLLMRequest struct {
	Command           string      `json:"command,omitempty"`
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
	AuthUUID          string      `json:"auth_uuid,omitempty"`
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

// PooledConnection wraps a WebSocket connection with metadata
type PooledConnection struct {
	Conn         *ws.Conn
	ConnectionID string
	CreatedAt    time.Time
	LastUsed     time.Time
	UseCount     int
	Identified   bool
}

// ConnectionPool manages multiple WebSocket connections with proper lifecycle
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
		maxUseCount:     100, // Recycle connection after 100 uses
		identifyTimeout: 5 * time.Second,
		ctx:             ctx,
		cancel:          cancel,
	}
	
	// Start cleanup goroutine
	go pool.cleanupStaleConnections()
	
	return pool
}

// Get retrieves a healthy connection from the pool
func (p *ConnectionPool) Get(authUUID string) (*PooledConnection, error) {
	if p.closed {
		return nil, fmt.Errorf("connection pool is closed")
	}

	// Try to get from pool first
	select {
	case pooledConn := <-p.pool:
		// Validate connection health
		if p.isHealthy(pooledConn) {
			pooledConn.LastUsed = time.Now()
			pooledConn.UseCount++
			
			p.mu.Lock()
			p.activeConns[pooledConn.ConnectionID] = pooledConn
			p.mu.Unlock()
			
			fmt.Printf("[Pool] Reusing connection %s (use #%d, age: %s)\n",
				pooledConn.ConnectionID[:8], pooledConn.UseCount,
				time.Since(pooledConn.CreatedAt).Round(time.Second))
			return pooledConn, nil
		}
		
		// Connection is unhealthy, close and create new one
		fmt.Printf("[Pool] Connection %s failed health check, discarding\n", 
			pooledConn.ConnectionID[:8])
		pooledConn.Conn.Close()
		
	default:
		// Pool is empty, check if we can create a new connection
	}

	// Create new connection
	p.mu.Lock()
	currentCount := len(p.activeConns)
	if currentCount >= p.maxConns {
		p.mu.Unlock()
		// Wait for available connection with timeout
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

// createConnection establishes a new WebSocket connection
func (p *ConnectionPool) createConnection(authUUID string) (*PooledConnection, error) {
	dialer := ws.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	
	conn, _, err := dialer.Dial(p.wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	// Generate connection ID
	connectionID := fmt.Sprintf("go-%d", time.Now().UnixNano())
	
	pooledConn := &PooledConnection{
		Conn:         conn,
		ConnectionID: connectionID,
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		UseCount:     0,
		Identified:   false,
	}

	// Identify connection to Python server
	if err := p.identifyConnection(pooledConn, authUUID); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to identify connection: %w", err)
	}

	p.mu.Lock()
	p.activeConns[connectionID] = pooledConn
	p.mu.Unlock()

	fmt.Printf("[Pool] Created new connection %s (total active: %d/%d)\n",
		connectionID[:8], len(p.activeConns), p.maxConns)

	return pooledConn, nil
}

// identifyConnection sends identification to Python server
func (p *ConnectionPool) identifyConnection(pooledConn *PooledConnection, authUUID string) error {
	identifyMsg := PythonLLMRequest{
		Command:  "identify",
		AuthUUID: authUUID,
	}

	// Set write deadline
	pooledConn.Conn.SetWriteDeadline(time.Now().Add(p.identifyTimeout))
	if err := pooledConn.Conn.WriteJSON(identifyMsg); err != nil {
		return fmt.Errorf("failed to send identify: %w", err)
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})

	// Wait for identification response
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
	fmt.Printf("[Pool] Connection %s identified successfully\n", pooledConn.ConnectionID[:8])
	
	return nil
}

// isHealthy checks if a connection is still usable
func (p *ConnectionPool) isHealthy(pooledConn *PooledConnection) bool {
	// Check age
	if time.Since(pooledConn.CreatedAt) > p.maxConnAge {
		fmt.Printf("[Pool] Connection %s exceeded max age\n", pooledConn.ConnectionID[:8])
		return false
	}

	// Check use count
	if pooledConn.UseCount >= p.maxUseCount {
		fmt.Printf("[Pool] Connection %s exceeded max use count\n", pooledConn.ConnectionID[:8])
		return false
	}

	// Ping check with short timeout
	deadline := time.Now().Add(2 * time.Second)
	if err := pooledConn.Conn.WriteControl(ws.PingMessage, []byte{}, deadline); err != nil {
		fmt.Printf("[Pool] Connection %s failed ping: %v\n", pooledConn.ConnectionID[:8], err)
		return false
	}

	return true
}

// Put returns a connection to the pool
func (p *ConnectionPool) Put(pooledConn *PooledConnection) {
	if pooledConn == nil || p.closed {
		return
	}

	p.mu.Lock()
	delete(p.activeConns, pooledConn.ConnectionID)
	p.mu.Unlock()

	// Validate before returning to pool
	if !p.isHealthy(pooledConn) {
		fmt.Printf("[Pool] Closing unhealthy connection %s instead of returning to pool\n",
			pooledConn.ConnectionID[:8])
		pooledConn.Conn.Close()
		return
	}

	select {
	case p.pool <- pooledConn:
		fmt.Printf("[Pool] Returned connection %s to pool (pool size: %d)\n",
			pooledConn.ConnectionID[:8], len(p.pool))
	default:
		// Pool is full, close the connection
		fmt.Printf("[Pool] Pool full, closing connection %s\n", pooledConn.ConnectionID[:8])
		pooledConn.Conn.Close()
	}
}

// Discard closes a connection without returning it to the pool
func (p *ConnectionPool) Discard(pooledConn *PooledConnection) {
	if pooledConn == nil {
		return
	}

	p.mu.Lock()
	delete(p.activeConns, pooledConn.ConnectionID)
	p.mu.Unlock()

	fmt.Printf("[Pool] Discarding connection %s\n", pooledConn.ConnectionID[:8])
	pooledConn.Conn.Close()
}

// cleanupStaleConnections periodically removes stale connections
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
				// Check if connection has been idle too long
				if now.Sub(conn.LastUsed) > 5*time.Minute {
					stale = append(stale, id)
				}
			}
			
			for _, id := range stale {
				conn := p.activeConns[id]
				delete(p.activeConns, id)
				conn.Conn.Close()
				fmt.Printf("[Pool] Cleaned up stale connection %s\n", id[:8])
			}
			p.mu.Unlock()

			// Also check pooled connections
			poolSize := len(p.pool)
			for i := 0; i < poolSize; i++ {
				select {
				case conn := <-p.pool:
					if !p.isHealthy(conn) {
						conn.Conn.Close()
						fmt.Printf("[Pool] Removed unhealthy pooled connection %s\n",
							conn.ConnectionID[:8])
					} else {
						// Return healthy connection
						select {
						case p.pool <- conn:
						default:
							conn.Conn.Close()
						}
					}
				default:
					// Pool empty
					break
				}
			}
		}
	}
}

// Close shuts down the connection pool
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	// Cancel cleanup goroutine
	p.cancel()

	// Close all active connections
	p.mu.Lock()
	for id, conn := range p.activeConns {
		conn.Conn.Close()
		delete(p.activeConns, id)
	}
	p.mu.Unlock()

	// Close pooled connections
	close(p.pool)
	for conn := range p.pool {
		conn.Conn.Close()
	}

	fmt.Printf("[Pool] Connection pool closed\n")
}

// GetStats returns pool statistics
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

// ChatService with improved connection pool
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
	fmt.Printf("[DEBUG] SendMessage started - ChatUUID: %s\n", data.ChatUUID)
	
	// Get connection from pool
	pooledConn, err := s.connPool.Get(authUUID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	// Track if connection should be returned or discarded
	shouldReturn := true
	defer func() {
		if shouldReturn {
			s.connPool.Put(pooledConn)
		} else {
			s.connPool.Discard(pooledConn)
		}
	}()

	chat := &d.Chat{ChatUUID: data.ChatUUID}
	err = s.r.GetByID(gctx, chat)
	if err != nil {
		return err
	}

	// Get agent configuration
	agent, err := s.agr.GetAgentByUUID(gctx, chat.AgentUUID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	// Set receiver info BEFORE saving
	data.ReceiverUUID = agent.AgentUUID
	data.ReceiverType = "AGENT"
	
	// Save user message to DB first
	if err := s.r.AttachMessage(gctx, data); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// Get chat history
	chat.History, err = s.r.GetChatHistory(gctx, chat.ChatUUID, s.lastMsgsLimit+1)
	if err != nil {
		return fmt.Errorf("failed to get chat history: %w", err)
	}
	
	fmt.Printf("[DEBUG] Chat history loaded: %d messages\n", len(chat.History))

	// Remove current message from history
	historyForPython := make([]d.Message, 0, len(chat.History))
	for _, msg := range chat.History {
		if msg.MessageUUID != data.MessageUUID {
			historyForPython = append(historyForPython, msg)
		}
	}
	
	fmt.Printf("[DEBUG] Sending %d history messages to Python\n", len(historyForPython))

	// Determine sync mode
	syncMode := s.determineChatHistoryStrategy(chat, uint64(len(historyForPython)))

	// Extract system prompt
	systemPrompt := "You are a helpful assistant."
	if agent.AgentConfig.AgentSystem.SystemPreset != nil {
		if prompt, ok := agent.AgentConfig.AgentSystem.SystemPreset["system_prompt"].(string); ok {
			systemPrompt = prompt
		}
	}

	// Build request for Python
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

	// Send request to Python with reasonable timeout
	fmt.Printf("[DEBUG] Sending request to Python...\n")
	pooledConn.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	if err := pooledConn.Conn.WriteJSON(request); err != nil {
		fmt.Printf("[DEBUG] Failed to send request: %v\n", err)
		shouldReturn = false // Connection is broken
		return fmt.Errorf("python service unavailable: %w", err)
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})
	fmt.Printf("[DEBUG] Request sent successfully\n")

	// Stream response chunks with extended timeout for LLM generation
	var agentMessageUUID string
	var messageContent d.MessageContent
	chunkCount := 0

	fmt.Printf("[DEBUG] Reading response chunks...\n")
	for {
		// 60 second timeout per chunk (LLMs can be slow)
		pooledConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		
		var response PythonLLMResponse
		err = pooledConn.Conn.ReadJSON(&response)
		if err != nil {
			fmt.Printf("[DEBUG] Failed to read after %d chunks: %v\n", chunkCount, err)
			shouldReturn = false // Connection likely broken
			return fmt.Errorf("failed to read response: %w", err)
		}
		
		chunkCount++

		if response.Error != "" {
			fmt.Printf("[DEBUG] Python error: %s\n", response.Error)
			// Don't discard connection on application error
			return fmt.Errorf("python service error: %s", response.Error)
		}

		if response.Partial {
			fmt.Printf("[DEBUG] Chunk #%d (len: %d)\n", chunkCount, len(response.Content))
			if streamCallback != nil {
				streamCallback(response.Content)
			}
		} else {
			fmt.Printf("[DEBUG] Final response after %d chunks\n", chunkCount)
			messageContent.Content = response.Content
			messageContent.MessageContentUUID = response.MessageContentUUID
			agentMessageUUID = response.MessageUUID
			break
		}
	}
	
	pooledConn.Conn.SetReadDeadline(time.Time{})
	fmt.Printf("[DEBUG] All chunks received\n")

	// Save agent response to DB
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
	fmt.Printf("[DEBUG] SendMessage completed successfully\n")
	return nil
}

func (s *ChatService) InitChat(gctx *gin.Context, data *d.Chat, streamCallback func(chunk string)) error {
	fmt.Printf("[DEBUG] InitChat started - ChatUUID: %s, AgentUUID: %s, AuthUUID: %s\n", 
		data.ChatUUID, data.AgentUUID, data.AuthUUID)

	if len(data.History) == 0 {
		return fmt.Errorf("at least one message is required to create a chat")
	}

	// Get connection from pool
	pooledConn, err := s.connPool.Get(data.AuthUUID)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to get connection: %v\n", err)
		return fmt.Errorf("failed to get connection: %w", err)
	}

	shouldReturn := true
	defer func() {
		if shouldReturn {
			s.connPool.Put(pooledConn)
		} else {
			s.connPool.Discard(pooledConn)
		}
	}()

	// Set timestamps
	now := time.Now()
	if data.CreatedAt.IsZero() {
		data.CreatedAt = now
	}
	if data.UpdatedAt.IsZero() {
		data.UpdatedAt = now
	}

	// Create chat in database
	fmt.Printf("[DEBUG] Creating chat in database...\n")
	if err := s.r.Create(gctx, data); err != nil {
		fmt.Printf("[DEBUG] Failed to create chat: %v\n", err)
		return fmt.Errorf("failed to create chat: %w", err)
	}
	fmt.Printf("[DEBUG] Chat created successfully\n")

	// Get agent configuration
	agent, err := s.agr.GetAgentByUUID(gctx, data.AgentUUID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	// Get first message
	userMessage := data.History[0]
	userMessage.ChatUUID = data.ChatUUID
	
	if userMessage.CreatedAt.IsZero() {
		userMessage.CreatedAt = now
	}

	// Save user message
	if err := s.r.AttachMessage(gctx, &userMessage); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// Extract system prompt
	systemPrompt := "You are a helpful assistant."
	if agent.AgentConfig.AgentSystem.SystemPreset != nil {
		if prompt, ok := agent.AgentConfig.AgentSystem.SystemPreset["system_prompt"].(string); ok {
			systemPrompt = prompt
		}
	}

	// Build request
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

	// Send request
	pooledConn.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	if err := pooledConn.Conn.WriteJSON(request); err != nil {
		shouldReturn = false
		return fmt.Errorf("python service unavailable: %w", err)
	}
	pooledConn.Conn.SetWriteDeadline(time.Time{})

	// Stream response
	var agentMessageUUID string
	var messageContent d.MessageContent
	chunkCount := 0
	
	for {
		pooledConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		
		var response PythonLLMResponse
		err = pooledConn.Conn.ReadJSON(&response)
		if err != nil {
			shouldReturn = false
			return fmt.Errorf("failed to read response: %w", err)
		}
		
		chunkCount++

		if response.Error != "" {
			return fmt.Errorf("python service error: %s", response.Error)
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

	// Save agent response
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
		return fmt.Errorf("failed to save agent message: %w", err)
	}

	data.History = []d.Message{userMessage, *agentMsg}
	
	fmt.Printf("[DEBUG] InitChat completed successfully\n")
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

// Cleanup should be called on service shutdown
func (s *ChatService) Cleanup() {
	s.connPool.Close()
}
