package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectName string
	Debug       bool
	SecretKey   string
	TimeZone    string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort string

	// MinIO
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	MinioBucket    string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))
	minioUseSSL, _ := strconv.ParseBool(getEnv("MINIO_USE_SSL", "false"))

	AppConfig = &Config{
		ProjectName: getEnv("PROJECT_NAME", "Pet Service API"),
		Debug:       debug,
		SecretKey:   getEnv("SECRET_KEY", "default-secret-key"),
		TimeZone:    getEnv("TIME_ZONE", "Asia/Ho_Chi_Minh"),

		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:     getEnv("POSTGRES_DB", "pet_service"),

		ServerPort: getEnv("SERVER_PORT", "8001"),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioUseSSL:    minioUseSSL,
		MinioBucket:    getEnv("MINIO_BUCKET", "pet-service"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
