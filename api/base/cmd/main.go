package main

import (
	db "aigents-base/internal/common/db"
	m "aigents-base/internal/auth-land/auth-signature/middleware"

	ah "aigents-base/internal/auth-land/auth/handlers"
	as "aigents-base/internal/auth-land/auth/services"
	ar "aigents-base/internal/auth-land/auth/repositories"

	"github.com/gin-contrib/cors"
	"time"
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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/auth")
	{
		auth.POST("/create", authHdlr.Create)
		auth.POST("/login", authHdlr.Login)
		auth.GET("/refresh", authHdlr.Refresh)
	}

	api := r.Group("/api/v1", m.AuthMiddleware())
	{
		agents := r.Group("/agents") {
			api.POST("/all", agentHdlr.Fetch)
			api.POST("/create", agentHdlr.Create)
			api.POST("/get", agentHdlr.GetByID)
		}

		chat := r.Group("/chat") {
			api.POST("/create", chatHdlr.Create)
			api.POST("/send-new-message", chatHdlr.SendMessage)
		}
	}

	r.Run(":8080")
}
