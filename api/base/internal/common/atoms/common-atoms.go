package atoms

import (
	"os"
	"log"
	"fmt"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Data      T      `json:"data,omitempty"`
}

func RespAtom[T any](
	gctx *gin.Context,
	code int,
	msg string,
	data T,
) {
	resp := Response[T]{
		Status:  code,
		Message: msg,
		Data: data,
	}

	gctx.JSON(code, resp)
}

func AbortAndBuildErrLogAtom(gctx *gin.Context, code int, abortMsg, logMsg string) error {
	err := BuildErrLogAtom(gctx, logMsg)
	AbortRespAtom(gctx, code, abortMsg)
	return err
}

func AbortRespAtom(gctx *gin.Context, code int, msg string) {
	gctx.AbortWithStatusJSON(code, map[string]any{
		"status": code,
		"error": msg,
	})
}

func BuildErrLogAtom(gctx *gin.Context, msg string) error {
	return fmt.Errorf(
		"[ERR LOG] %s | IP: %s | UA: %s | P: %s -> %s",
		time.Now().Format("2006-01-02 15:04:05"),
		gctx.ClientIP(),
		gctx.Request.UserAgent(),
		gctx.FullPath(),
	msg)
}

func ParseEnvMinutesAtom(eVar string, fallback int) time.Duration {
	valStr := os.Getenv(eVar)
	if valStr == "" {
		return time.Duration(fallback) * time.Minute
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid minutes value for %s: %v", eVar, err)
	}

	return time.Duration(val) * time.Minute
}

func FeedErrLogToFile(err error) {
	if err == nil {
		return
	}

	spath := os.Getenv("ERR_LOG_FPATH")
	if spath == "" {
		spath = "ERR_LOG"
	}

	f, openErr := os.OpenFile(spath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		log.Fatalf("failed to open error log file: %w", openErr)
		return
	}
	defer f.Close()

	_, writeErr := f.WriteString(err.Error() + "\n")
	if writeErr != nil {
		log.Fatalf("failed to write to error log file: %w", writeErr)
		return
	}
}
