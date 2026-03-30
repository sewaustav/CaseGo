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

	PublicKey *rsa.PublicKey
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	publicKeyStr := os.Getenv("PUBLIC_KEY")
	if publicKeyStr == "" {
		panic("PUBLIC_KEY environment variable not set")
	}

	publicKey, err := ParseRSAPublicKey(publicKeyStr)
	if err != nil {
		log.Fatal(publicKeyStr)
	}

	dbPortStr := os.Getenv("POSTGRES_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Failed to parse DB_PORT")
	}

	return &Config{
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBPort:     dbPort,
		PublicKey:  publicKey,
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
