package jwt

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type JwtService interface {
	GenerateToken(userID int64, role models.UserRole) (string, error)
}

type Token struct {
	privateKey *rsa.PrivateKey
}

func NewToken(privateKey *rsa.PrivateKey) *Token {
	return &Token{privateKey: privateKey}
}

func (t *Token) GenerateToken(userID int64, role models.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"iss":  "cases",                       // Издатель
		"aud":  "profile",                     // Аудитория
		"sub":  strconv.FormatInt(userID, 10), // ID пользователя как string
		"role": string(role),                  // Роль как string
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Теперь просто передаем объект ключа, парсинг больше не нужен
	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
