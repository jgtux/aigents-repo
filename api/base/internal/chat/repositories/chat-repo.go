package repositories

import (
	d "aigents-base/internal/chat/domain"
	chitf "aigents-base/internal/chat/interfaces"
	c_at "aigents-base/internal/common/atoms"
	"database/sql"
	"fmt"
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
	var agentExists, authExists bool

	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM agents WHERE agent_uuid = $1 AND deleted_at IS NULL)", data.AgentUUID).Scan(&agentExists)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not verify agent existence. Failed to check agent existence: %s", err.Error()))
		return err
	}

	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM auths WHERE auth_uuid = $1 AND deleted_at IS NULL)", data.AuthUUID).Scan(&authExists)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not verify auth existence. Failed to check auth existence: %s", err.Error()))
		return err
	}

	if !agentExists {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Agent not found. Agent with UUID %s does not exist", data.AgentUUID))
		return err
	}

	if !authExists {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Authentication not found. Auth with UUID %s does not exist", data.AuthUUID))
		return err
	}

	query := `
		INSERT INTO chats (chat_uuid, agent_uuid, auth_uuid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_uuid) DO NOTHING
		RETURNING chat_uuid
	`

	var returnedUUID string
	err = r.db.QueryRow(query,
		data.ChatUUID,
		data.AgentUUID,
		data.AuthUUID,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&returnedUUID)

	if err != nil {
		if err == sql.ErrNoRows {
			var existingChatUUID string
			err = r.db.QueryRow("SELECT chat_uuid FROM chats WHERE chat_uuid = $1", data.ChatUUID).Scan(&existingChatUUID)
			if err != nil {
				err = c_at.BuildErrLogAtom(
					gctx,
					fmt.Sprintf("(R) Could not create chat. Chat conflict but doesn't exist: %s", err.Error()))
				return err
			}
			return nil
		}

		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not create chat. Failed to create chat: %s", err.Error()))
		return err
	}

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

	if err == sql.ErrNoRows {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Chat not found. Chat with UUID %s not found", data.ChatUUID))
		return err
	}

	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not get chat. Failed to get chat: %s", err.Error()))
		return err
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
	var chatExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chats WHERE chat_uuid = $1)", msg.ChatUUID).Scan(&chatExists)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not verify chat existence. Failed to check chat existence: %s", err.Error()))
		return err
	}

	if !chatExists {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Chat not found. Chat with UUID %s does not exist", msg.ChatUUID))
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not begin transaction. Failed to begin transaction: %s", err.Error()))
		return err
	}
	defer tx.Rollback()

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

	_, err = tx.Exec(insertSQL,
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
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not attach message. Failed to insert message: %s", err.Error()))
		return err
	}

	updateSQL := `
		UPDATE chats
		SET updated_at = NOW()
		WHERE chat_uuid = $1
	`

	_, err = tx.Exec(updateSQL, msg.ChatUUID)
	if err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not update chat timestamp. Failed to update chat timestamp: %s", err.Error()))
		return err
	}

	if err = tx.Commit(); err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("(R) Could not commit transaction. Failed to commit transaction: %s", err.Error()))
		return err
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
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not fetch chat history. Failed to query chat history: %s", err.Error()))
		return nil, err
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
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("Could not fetch chat history. Failed to scan message: %s", err.Error()))
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	if err = rows.Err(); err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not fetch chat history. Row iteration failed: %s", err.Error()))
		return nil, err
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
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not fetch recent messages. Failed to query recent messages: %s", err.Error()))
		return nil, err
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
			err = c_at.BuildErrLogAtom(
				gctx,
				fmt.Sprintf("Could not fetch recent messages. Failed to scan message: %s", err.Error()))
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	if err = rows.Err(); err != nil {
		err = c_at.BuildErrLogAtom(
			gctx,
			fmt.Sprintf("Could not fetch recent messages. Row iteration failed: %s", err.Error()))
		return nil, err
	}

	return msgs, nil
}
