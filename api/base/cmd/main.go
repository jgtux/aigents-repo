package main

import (
	m "aigents-base/internal/auth-land/auth-signature/middleware"
	db "aigents-base/internal/common/db"
	"os"

	ah "aigents-base/internal/auth-land/auth/handlers"
	ar "aigents-base/internal/auth-land/auth/repositories"
	as "aigents-base/internal/auth-land/auth/services"

	agh "aigents-base/internal/agents/handlers"
	agr "aigents-base/internal/agents/repositories"
	ags "aigents-base/internal/agents/services"

	chh "aigents-base/internal/chat/handlers"
	chr "aigents-base/internal/chat/repositories"
	chs "aigents-base/internal/chat/services"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func main() {
	db.Init()

	authRepo := ar.NewAuthRepository(db.DB)
	authSv := as.NewAuthService(authRepo)
	authHdlr := ah.NewAuthHandler(authSv)

	agentRepo := agr.NewAgentRepository(db.DB)
	agentSv := ags.NewAgentService(agentRepo)
	agentHdlr := agh.NewAgentHandler(agentSv)

	chatRepo := chr.NewChatRepository(db.DB)
	chatSv := chs.NewChatService(chatRepo, agentRepo, os.Getenv("WS_AI_MS_URL"), 20, 20)
	chatHdlr := chh.NewChatHandler(chatSv)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/create", authHdlr.Create)
		auth.POST("/login", authHdlr.Login)
		auth.GET("/refresh", authHdlr.Refresh)
	}

	api := r.Group("/api/v1", m.AuthMiddleware())
	{
		agents := api.Group("/agents")
		{
			agents.POST("/all", agentHdlr.Fetch)
			agents.POST("/create", agentHdlr.Create)
			agents.POST("/get", agentHdlr.GetByID)
		}

		chat := api.Group("/chat")
		{
			chat.POST("/create", chatHdlr.Create)
			chat.POST("/send-new-message", chatHdlr.SendMessage)
		}
	}

	r.Run(":8080")
}
