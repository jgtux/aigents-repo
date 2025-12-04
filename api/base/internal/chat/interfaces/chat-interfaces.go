package interfaces

import (
	"time"
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/chat/domain"
	"github.com/gin-gonic/gin"
)

type ChatServiceITF interface {
	citf.Common[d.Chat]
	SendMessage(gctx *gin.Context, data *d.Message, authUUID string, streamCallback func(chunk string)) error
	InitChat(gctx *gin.Context, data *d.Chat, streamCallback func(chunk string)) error
}

type ChatRepositoryITF interface {
	citf.Common[d.Chat]
	AttachMessage(gctx *gin.Context, msg *d.Message) error
	GetChatHistory(gctx *gin.Context, chatUUID string, limit uint64) ([]d.Message, error)
	GetRecentMessages(gctx *gin.Context, chatUUID string, since time.Time, limit uint64) ([]d.Message, error)
}
