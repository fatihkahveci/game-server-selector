package services

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigService interface {
	GetHttpPort() string
	GetStorageType() string
	GetSslCertPath() string
	GetSslKeyPath() string
	GetServerDirtySeconds() int
}

type configService struct{}

func NewConfigService() ConfigService {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &configService{}
}

func (c *configService) GetHttpPort() string {
	return ":" + c.getEnv("HTTP_PORT", "80")
}

func (c *configService) GetStorageType() string {
	return c.getEnv("STORAGE_TYPE", "inmemory")
}

func (c *configService) GetSslCertPath() string {
	return c.getEnv("SSL_CERT", "")
}

func (c *configService) GetSslKeyPath() string {
	return c.getEnv("SSL_KEY", "")
}

func (c *configService) GetServerDirtySeconds() int {
	return c.getEnvInt("SERVER_DIRTY_SECONDS", 60)
}

func (c *configService) getEnv(key string, fallback string) string {
	env := os.Getenv(key)
	if env != "" {
		return env
	}
	return fallback
}

func (c *configService) getEnvInt(key string, fallback int) int {
	env := os.Getenv(key)
	if env != "" {
		intEnv, err := strconv.Atoi(env)
		if err != nil {
			log.Default().Println("Error converting env to int", err)
			return fallback
		}
		return intEnv
	}
	return fallback
}
