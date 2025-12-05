package repositories

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	"fmt"

	"database/sql"

	"time"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) chitf.ChatRepositoryITF {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(gctx *gin.Context, data *d.Chat) error {
	fmt.Printf("[DEBUG REPO] Creating chat - UUID: %s, AgentUUID: %s, AuthUUID: %s, CreatedAt: %v, UpdatedAt: %v\n",
		data.ChatUUID, data.AgentUUID, data.AuthUUID, data.CreatedAt, data.UpdatedAt)
	
	// Primeiro verifica se o agent e auth existem
	var agentExists, authExists bool
	
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM agents WHERE agent_uuid = $1 AND deleted_at IS NULL)", data.AgentUUID).Scan(&agentExists)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Error checking agent existence: %v\n", err)
		return fmt.Errorf("failed to check agent existence: %w", err)
	}
	fmt.Printf("[DEBUG REPO] Agent exists: %v\n", agentExists)
	
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM auths WHERE auth_uuid = $1 AND deleted_at IS NULL)", data.AuthUUID).Scan(&authExists)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Error checking auth existence: %v\n", err)
		return fmt.Errorf("failed to check auth existence: %w", err)
	}
	fmt.Printf("[DEBUG REPO] Auth exists: %v\n", authExists)
	
	if !agentExists {
		return fmt.Errorf("agent with UUID %s does not exist", data.AgentUUID)
	}
	if !authExists {
		return fmt.Errorf("auth with UUID %s does not exist", data.AuthUUID)
	}
	
	// Tenta inserir o chat
	query := `
		INSERT INTO chats (chat_uuid, agent_uuid, auth_uuid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_uuid) DO NOTHING
		RETURNING chat_uuid
	`
	var returnedUUID string
	err = r.db.QueryRowContext(gctx, query,
		data.ChatUUID,
		data.AgentUUID,
		data.AuthUUID,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&returnedUUID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[DEBUG REPO] Chat already exists (conflict), checking if it exists in DB...\n")
			// Verifica se o chat realmente existe
			var existingChatUUID string
			err = r.db.QueryRow("SELECT chat_uuid FROM chats WHERE chat_uuid = $1", data.ChatUUID).Scan(&existingChatUUID)
			if err != nil {
				fmt.Printf("[DEBUG REPO] Chat doesn't exist even after conflict: %v\n", err)
				return fmt.Errorf("chat conflict but doesn't exist: %w", err)
			}
			fmt.Printf("[DEBUG REPO] Chat already exists with UUID: %s\n", existingChatUUID)
			return nil
		}
		fmt.Printf("[DEBUG REPO] Failed to create chat: %v\n", err)
		return fmt.Errorf("failed to create chat: %w", err)
	}
	
	fmt.Printf("[DEBUG REPO] Chat created successfully with UUID: %s\n", returnedUUID)
	return nil
}

func (r *ChatRepository) GetByID(gctx *gin.Context, data *d.Chat) error {
	fmt.Printf("[DEBUG REPO] GetByID - ChatUUID: %s\n", data.ChatUUID)
	
	query := `
		SELECT chat_uuid, agent_uuid, auth_uuid, created_at, updated_at
		FROM chats
		WHERE chat_uuid = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, data.ChatUUID).Scan(
		&data.ChatUUID,
		&data.AgentUUID,
		&data.AuthUUID,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("[DEBUG REPO] GetByID failed: %v\n", err)
		return fmt.Errorf("failed to get chat: %w", err)
	}
	
	fmt.Printf("[DEBUG REPO] GetByID success - AgentUUID: %s, AuthUUID: %s\n", data.AgentUUID, data.AuthUUID)
	return nil
}

func (r *ChatRepository) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Chat, error) {
	return nil, nil
}

func (r *ChatRepository) Update(gctx *gin.Context, data *d.Chat) error {
	return nil
}

func (r *ChatRepository) Delete(gctx *gin.Context, data *d.Chat) error {
	return nil
}

func (r *ChatRepository) AttachMessage(gctx *gin.Context, msg *d.Message) error {
	fmt.Printf("[DEBUG REPO] AttachMessage - MessageUUID: %s, ChatUUID: %s, Sender: %s (%s), Receiver: %s (%s)\n",
		msg.MessageUUID, msg.ChatUUID, msg.SenderUUID, msg.SenderType, msg.ReceiverUUID, msg.ReceiverType)
	fmt.Printf("[DEBUG REPO] Message content UUID: %s, Content length: %d\n",
		msg.MessageContent.MessageContentUUID, len(msg.MessageContent.Content))
	
	// Verifica se o chat existe antes de inserir a mensagem
	var chatExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chats WHERE chat_uuid = $1)", msg.ChatUUID).Scan(&chatExists)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Error checking chat existence: %v\n", err)
		return fmt.Errorf("failed to check chat existence: %w", err)
	}
	fmt.Printf("[DEBUG REPO] Chat exists before insert: %v\n", chatExists)
	
	if !chatExists {
		return fmt.Errorf("chat with UUID %s does not exist", msg.ChatUUID)
	}
	
	// Inicia a transação
	tx, err := r.db.BeginTx(gctx, nil)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Failed to begin transaction: %v\n", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert message content e message em uma única query
	insertSQL := `
		WITH inserted_content AS (
			INSERT INTO message_contents (message_content_uuid, message_content)
			VALUES ($1, $2)
			RETURNING message_content_uuid
		)
		INSERT INTO messages (
			message_uuid, sender_uuid, sender_type, receiver_uuid, receiver_type,
			chat_uuid, message_content_uuid, created_at
		)
		SELECT $3, $4, $5, $6, $7, $8, message_content_uuid, $9
		FROM inserted_content
	`

	_, err = tx.ExecContext(gctx, insertSQL,
		msg.MessageContent.MessageContentUUID,
		msg.MessageContent.Content,
		msg.MessageUUID,
		msg.SenderUUID,
		msg.SenderType,
		msg.ReceiverUUID,
		msg.ReceiverType,
		msg.ChatUUID,
		msg.CreatedAt,
	)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Failed to insert message: %v\n", err)
		return fmt.Errorf("failed to insert message: %w", err)
	}
	fmt.Printf("[DEBUG REPO] Message inserted successfully\n")

	// Atualiza o timestamp do chat
	updateSQL := `
		UPDATE chats
		SET updated_at = NOW()
		WHERE chat_uuid = $1
	`
	result, err := tx.ExecContext(gctx, updateSQL, msg.ChatUUID)
	if err != nil {
		fmt.Printf("[DEBUG REPO] Failed to update chat timestamp: %v\n", err)
		return fmt.Errorf("failed to update chat timestamp: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("[DEBUG REPO] Chat updated, rows affected: %d\n", rowsAffected)

	// Commit da transação
	if err = tx.Commit(); err != nil {
		fmt.Printf("[DEBUG REPO] Failed to commit transaction: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("[DEBUG REPO] AttachMessage completed successfully\n")
	return nil
}

func (r *ChatRepository) GetChatHistory(gctx *gin.Context, chatUUID string, limit uint64) ([]d.Message, error) {
	// FIXED: Get the LAST N messages in correct chronological order
	query := `
		SELECT 
			m.message_uuid,
			m.sender_uuid,
			m.sender_type,
			m.receiver_uuid,
			m.receiver_type,
			m.chat_uuid,
			m.message_content_uuid,
			mc.message_content,
			m.created_at
		FROM (
			SELECT *
			FROM messages
			WHERE chat_uuid = $1
			ORDER BY created_at DESC
			LIMIT $2
		) m
		INNER JOIN message_contents mc ON m.message_content_uuid = mc.message_content_uuid
		ORDER BY m.created_at ASC
	`

	rows, err := r.db.Query(query, chatUUID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat history: %w", err)
	}
	defer rows.Close()

	var msgs []d.Message
	for rows.Next() {
		var msg d.Message
		err := rows.Scan(
			&msg.MessageUUID,
			&msg.SenderUUID,
			&msg.SenderType,
			&msg.ReceiverUUID,
			&msg.ReceiverType,
			&msg.ChatUUID,
			&msg.MessageContent.MessageContentUUID,
			&msg.MessageContent.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}


func (r *ChatRepository) GetRecentMessages(gctx *gin.Context, chatUUID string, since time.Time, limit uint64) ([]d.Message, error) {
	query := `
		SELECT 
			m.message_uuid,
			m.sender_uuid,
			m.sender_type,
			m.receiver_uuid,
			m.receiver_type,
			m.chat_uuid,
			m.message_content_uuid,
			mc.message_content,
			m.created_at
		FROM messages m
		INNER JOIN message_contents mc ON m.message_content_uuid = mc.message_content_uuid
		WHERE m.chat_uuid = $1 AND m.created_at > $2
		ORDER BY m.created_at ASC
		LIMIT $3
	`

	rows, err := r.db.Query(query, chatUUID, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
	}
	defer rows.Close()

	var msgs []d.Message
	for rows.Next() {
		var msg d.Message
		err := rows.Scan(
			&msg.MessageUUID,
			&msg.SenderUUID,
			&msg.SenderType,
			&msg.ReceiverUUID,
			&msg.ReceiverType,
			&msg.ChatUUID,
			&msg.MessageContent.MessageContentUUID,
			&msg.MessageContent.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
