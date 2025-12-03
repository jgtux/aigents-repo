package middleware

import (
	c_at "aigents-base/internal/common/atoms"

	"os"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UUID string `json:"auth_uuid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	RefreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
	AccessTokenTTL  = c_at.ParseEnvMinutesAtom("ACCESS_TOKEN_TTL", 15)
	RefreshTokenTTL = c_at.ParseEnvMinutesAtom("REFRESH_TOKEN_TTL", 10080)
)

func AuthMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		tokenStr, err := gctx.Cookie("access_token")
		if err != nil {
			err := c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusUnauthorized,
				"(M) Missing token.",
				"Missing token.")
			c_at.FeedErrLogToFile(err)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			err := c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusUnauthorized,
				"(M) Invalid token.",
				"Invalid token.")
			c_at.FeedErrLogToFile(err)
			return
		}
		gctx.Set("auth_uuid", claims.UUID)
		gctx.Set("role", claims.Role)

		gctx.Next()
	}
}

func AuthorizeRole(allowedRoles map[string]bool) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		role, exists := gctx.Get("role")
		if !exists {
			err := c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusForbidden,
				"(M) Insufficient role.",
				"Insufficient role.")
			c_at.FeedErrLogToFile(err)
			return
		}

		roleStr, ok := role.(string)
		if !ok || !allowedRoles[roleStr] {
			err := c_at.AbortAndBuildErrLogAtom(
				gctx,
				http.StatusForbidden,
				"(M) Insufficient role.",
				"Insufficient role.")
			c_at.FeedErrLogToFile(err)
			return
		}

		gctx.Next()
	}
}

func GenerateJWT(gctx *gin.Context, c *Claims, useRefresh bool) (string, error) {
	now := time.Now()

	var ttl time.Duration
	var secret []byte

	if useRefresh {
		ttl = RefreshTokenTTL
		secret = RefreshSecret
	} else {
		ttl = AccessTokenTTL
		secret = JWTSecret
	}

	c.IssuedAt = jwt.NewNumericDate(now)
	c.ExpiresAt = jwt.NewNumericDate(now.Add(ttl))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedStr, err := token.SignedString(secret)
	if err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(M) Could not generate token.",
			fmt.Sprintf("Could not generate token: %s", err.Error()))
		return "", err
	}

	return signedStr, nil
}

func GetAuthUUID(gctx *gin.Context) (string, bool) {
	val, exists := gctx.Get("auth_uuid")
	if !exists {
		return "", false
	}

	uuidStr, ok := val.(string)
	if !ok {
		return "", false
	}

	return uuidStr, true
}
