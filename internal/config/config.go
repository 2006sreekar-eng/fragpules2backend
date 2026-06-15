package config

import (
	"os"
	"path/filepath"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	CORSAllowedOrigins string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	} else if _, err := os.Stat(filepath.Join("..", ".env")); err == nil {
		_ = godotenv.Load(filepath.Join("..", ".env"))
	}

	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "fragpulse_user"),
		DBPassword:         getEnv("DB_PASSWORD", "fragpulse_password"),
		DBName:             getEnv("DB_NAME", "fragpulse_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}