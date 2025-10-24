package handlers

import (
	d "aigents-base/internal/auth-land/auth/domain"
	auitf "aigents-base/internal/auth-land/auth/interfaces"
	c_at "aigents-base/internal/common/atoms"
	m "aigents-base/internal/auth-land/auth-signature/middleware"

	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	s auitf.AuthServiceITF
}

func NewAuthHandler(sv auitf.AuthServiceITF) *AuthHandler {
	return &AuthHandler{s: sv}
}

func (h *AuthHandler) Create(gctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		c_at.AbortRespAtom(gctx, http.StatusBadRequest, "(H) Invalid body request.")
		return
	}

	err := h.s.Create(&d.Auth{Email: req.Email, Password: req.Password})
	if err != nil {
		err(gctx)
		return
	}

	c_at.RespAtom(gctx, http.StatusCreated, "(*) Authentication created.")
}

func (h *AuthHandler) Login(gctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := gctx.ShouldBindJSON(&req); err != nil {
		c_at.AbortRespAtom(gctx, http.StatusBadRequest, "(H) Invalid body request.")
		return
	}

	auth := &d.Auth{Email: req.Email, Password: req.Password}
	err := h.s.Comparate(auth)
	if err != nil {
		err(gctx)
		return
	}

	claims := &m.Claims{ UUID: auth.UUID, Role: auth.Role }
	accessToken, err := m.GenerateJWT(claims, false)
	if err != nil {
		err(gctx)
		return
	}

	refreshToken, err := m.GenerateJWT(claims, true)
	if err != nil {
		err(gctx)
		return
	}

	gctx.SetCookie("access_token", accessToken, int(m.AccessTokenTTL.Seconds()), "/", "", false, true)
	gctx.SetCookie("refresh_token", refreshToken, int(m.RefreshTokenTTL.Seconds()), "/", "", false, true)

	c_at.RespAtom(gctx, http.StatusOK, "(*) Login successful.")
}

func (h *AuthHandler) Refresh(gctx *gin.Context) {
	refreshToken, err := gctx.Cookie("refresh_token")
	if err != nil {
		c_at.AbortRespAtom(gctx, http.StatusUnauthorized, "(H) Missing refresh token.")
		return
	}

	claims := &m.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		return m.RefreshSecret, nil
	})

	if err != nil || !token.Valid {
		c_at.AbortRespAtom(gctx, http.StatusUnauthorized, "(M) Invalid refresh token.")
		return
	}

	newAccessToken, errF := m.GenerateJWT(claims, false)
	if errF != nil {
		errF(gctx)
		return
	}

	gctx.SetCookie("access_token", newAccessToken, int(m.AccessTokenTTL.Seconds()), "/", "", false, true)

	c_at.RespAtom(gctx, http.StatusOK, "(*) Access token refreshed.")
}

func (h *AuthHandler) GetByID(gctx *gin.Context) {
}

func (h *AuthHandler) Fetch(gctx *gin.Context) {
}

func (h *AuthHandler) Update(gctx *gin.Context, data *d.Auth) {
}

func (h *AuthHandler) Delete(gctx *gin.Context) {
}
