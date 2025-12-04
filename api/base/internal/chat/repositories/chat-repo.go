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
	return nil
}


func (r *ChatRepository) GetByID(gctx *gin.Context, data *d.Chat) error {
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
		return nil
	}
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
	sql := `
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
		FROM inserted_content;

		UPDATE chats
		SET updated_at = NOW()
		WHERE chat_uuid = $8;
	`

	_, err := r.db.ExecContext(gctx, sql,
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
		return fmt.Errorf("failed to insert message and update chat: %w", err)
	}

	return nil
}



func (r *ChatRepository) GetChatHistory(gctx *gin.Context, chatUUID string, limit uint64) ([]d.Message, error) {
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
		WHERE m.chat_uuid = $1
		ORDER BY m.created_at ASC
		LIMIT $2
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
