package interfaces

import (
	"time"
	citf "aigents-base/internal/common/interfaces"
	d "aigents-base/internal/chat/domain"
	"github.com/gin-gonic/gin"
)

type ChatServiceITF interface {
	citf.Common[d.Chat]
}

type ChatRepositoryITF interface {
	citf.Common[d.Chat]
	AttachMessage(gctx *gin.Context, msg *d.Message) error
	GetChatHistory(gctx *gin.Context, data *d.Chat, limit uint64) error
	GetRecentMessages(gctx *gin.Context, data *d.Chat, since time.Time, limit uint64) error
}
