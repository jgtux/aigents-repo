package main

import (
	db "aigents-base/internal/common/db"

	ah "aigents-base/internal/auth-land/auth/handlers"
	as "aigents-base/internal/auth-land/auth/services"
	ar "aigents-base/internal/auth-land/auth/repositories"

	"github.com/gin-gonic/gin"
)


func main() {

	authRepo := ar.NewAuthRepository(db.DB)
	authSv := as.NewAuthService(authRepo)
	authHdlr := ah.NewAuthHandler(authSv)


	r := gin.Default()

	api := r.Group("/api/v1")

	{
		api.POST("/create", authHdlr.Create)
		api.POST("/login", authHdlr.Login)
		api.GET("/refresh", authHdlr.Refresh)
	}

	r.Run(":8080")
}
