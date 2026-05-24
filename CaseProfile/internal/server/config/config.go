package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int

	// AuthPublicKey — ключ Auth-сервиса, которым подписаны пользовательские JWT.
	// Используется HTTP-middleware для проверки запросов от Flutter-клиента.
	AuthPublicKey *rsa.PublicKey

	// CaseGoPublicKey — ключ CaseGo, которым подписаны сервисные JWT.
	// Используется gRPC-interceptor для проверки запросов от CaseGo.
	CaseGoPublicKey *rsa.PublicKey
}

func LoadConfig() *Config {
	godotenv.Load()

	// Auth public key (для HTTP-роутов)
	authKeyStr := os.Getenv("AUTH_PUBLIC_KEY")
	if authKeyStr == "" {
		panic("AUTH_PUBLIC_KEY environment variable not set")
	}
	authPublicKey, err := ParseRSAPublicKey(authKeyStr)
	if err != nil {
		log.Fatalf("Failed to parse AUTH_PUBLIC_KEY: %v", err)
	}

	// CaseGo public key (для gRPC)
	caseGoKeyStr := os.Getenv("CASEGO_PUBLIC_KEY")
	if caseGoKeyStr == "" {
		panic("CASEGO_PUBLIC_KEY environment variable not set")
	}
	caseGoPublicKey, err := ParseRSAPublicKey(caseGoKeyStr)
	if err != nil {
		log.Fatalf("Failed to parse CASEGO_PUBLIC_KEY: %v", err)
	}

	dbPortStr := os.Getenv("POSTGRES_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Failed to parse DB_PORT")
	}

	return &Config{
		DBHost:          os.Getenv("POSTGRES_HOST"),
		DBName:          os.Getenv("POSTGRES_DB"),
		DBUser:          os.Getenv("POSTGRES_USER"),
		DBPassword:      os.Getenv("POSTGRES_PASSWORD"),
		DBPort:          dbPort,
		AuthPublicKey:   authPublicKey,
		CaseGoPublicKey: caseGoPublicKey,
	}
}

func ParseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unknown type of public key")
	}
}
