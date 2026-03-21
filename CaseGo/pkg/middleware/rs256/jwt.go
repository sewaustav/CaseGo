package rs256

import (
	"crypto/rsa"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "sub"
	RoleKey   = "role"
)

type JWTAuthMiddleware struct {
	publicKey *rsa.PublicKey
	issuer    string
	audience  string
	logger    *slog.Logger
}

func New(pubKey *rsa.PublicKey, issuer, audience string) *JWTAuthMiddleware {
	logger := slog.Default()

	return &JWTAuthMiddleware{
		publicKey: pubKey,
		issuer:    issuer,
		audience:  audience,
		logger:    logger.With("component", "jwt-auth"),
	}
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Debug("missing authorization header", "path", c.Request.URL.Path)
			unauthorized(c)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			m.logger.Warn("invalid auth header format", "header", authHeader)
			unauthorized(c)
			return
		}

		claims, err := m.verifyToken(parts[1])
		if err != nil {
			// Логируем ошибку проверки. Сюда попадут истекшие токены, кривые подписи и т.д.
			m.logger.Info("token verification failed", "err", err, "client_ip", c.ClientIP())
			unauthorized(c)
			return
		}

		// Логируем успешный вход на уровне Debug, чтобы не засорять продакшн-логи
		m.logger.Debug("user authenticated", "user_id", claims.UserID, "role", claims.Role)

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}

type tokenClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"user_role"`
	jwt.RegisteredClaims
}

type Claims struct {
	UserID int64
	Role   int
}

func (m *JWTAuthMiddleware) verifyToken(tokenStr string) (*Claims, error) {
	if m.issuer == "" || m.audience == "" {
		m.logger.Error("middleware configuration error: missing issuer or audience")
		return nil, errors.New("auth middleware is not properly configured")
	}

	tokenClaims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		tokenClaims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.publicKey, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	id, err := strconv.ParseInt(tokenClaims.UserID, 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse int")
	}
	if id <= 0 {
		return nil, errors.New("invalid user id in token")
	}

	role, err := strconv.Atoi(tokenClaims.Role)
	if err != nil {
		return nil, errors.New("failed to parse role: " + err.Error())
	}
	if role < 0 || role > 3 {
		return nil, errors.New("invalid role in token")
	}

	return &Claims{
		UserID: id,
		Role:   role,
	}, nil
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{"error": "unauthorized"},
	)
}
