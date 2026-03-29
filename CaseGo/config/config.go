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

	RedisHost     string
	RedisPort     int
	RedisPassword string

	LLMURL string

	GRPCSEVER string

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadConfig() *Config {
	godotenv.Load()

	publicKeyStr := os.Getenv("PUBLIC_KEY")
	if publicKeyStr == "" {
		panic("PUBLIC_KEY environment variable not set")
	}

	publicKey, err := ParseRSAPublicKey(publicKeyStr)
	if err != nil {
		log.Fatal(publicKeyStr)
	}

	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		panic("PRIVATE_KEY environment variable not set")
	}

	privateKey, err := ParseRSAPrivateKey(privateKeyStr)
	if err != nil {
		log.Fatal("Failed to parse PRIVATE_KEY")
	}

	dbPortStr := os.Getenv("POSTGRES_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Failed to parse DB_PORT")
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("Failed to parse REDIS_PORT")
	}

	return &Config{
		DBHost:        os.Getenv("POSTGRES_HOST"),
		DBName:        os.Getenv("POSTGRES_DB"),
		DBUser:        os.Getenv("POSTGRES_USER"),
		DBPassword:    os.Getenv("POSTGRES_PASSWORD"),
		DBPort:        dbPort,
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisPort,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		LLMURL:        os.Getenv("LLM_URL"),
		GRPCSEVER:     os.Getenv("GRPC_SEVER"),
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

func ParseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch key := key.(type) {
	case *rsa.PrivateKey:
		return key, nil
	default:
		return nil, errors.New("unknown type of private key")
	}
}
