package atoms

import (
	"os"
	"log"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
)

func RespAtom(gctx *gin.Context, code int, msg string) {
	gctx.JSON(code, map[string]any{
		"status": code,
		"message": msg,
	})
}

func AbortRespAtom(gctx *gin.Context, code int, msg string) {
	gctx.AbortWithStatusJSON(code, map[string]any{
		"status": code,
		"error": msg,
	})
}

func RespFuncAbortAtom(code int, msg string) func(*gin.Context) {
	return func(gctx *gin.Context) {
		gctx.AbortWithStatusJSON(code, map[string]any{
			"status": code,
			"error": msg,
		})
	}
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
