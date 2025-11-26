package main

import (
	db "aigents-base/internal/common/db"

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


	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // or {"*"} for dev only
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you need cookies/auth
		MaxAge:           12 * time.Hour,
	}))
	

	api := r.Group("/api/v1")

	{
		api.POST("/auth/create", authHdlr.Create)
		api.POST("/auth/login", authHdlr.Login)
		api.GET("/auth/refresh", authHdlr.Refresh)
	}

	r.Run(":8080")
}
