package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessToken  string
	SecureKey    string
	UPLOAD_PATH  string
	DB_HOST      string
	DB_PORT      string
	DB_USER      string
	DB_PASSWORD  string
	DB_NAME      string
	DB_MAX_CONNS int

	MINIO_ENDPOINT   string
	MINIO_ACCESS_KEY string
	MINIO_SECRET_KEY string
	MINIO_SECURE     bool
	MINIO_BUCKET     string
	MINIO_WORKERS    int
}

func Init() *Config {
	var cfg Config
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.AccessToken = loadEnvVariable("ACCESS_TOKEN")
	cfg.DB_HOST = loadEnvVariable("DB_HOST")
	cfg.DB_PORT = loadEnvVariable("DB_PORT")
	cfg.DB_USER = loadEnvVariable("DB_USER")
	cfg.DB_PASSWORD = loadEnvVariable("DB_PASSWORD")
	cfg.DB_NAME = loadEnvVariable("DB_NAME")
	cfg.DB_MAX_CONNS = loadEnvVariableInt("DB_MAX_CONNS")
	cfg.SecureKey = "w3r1Sec4re_Token_"

	cfg.MINIO_ENDPOINT = loadEnvVariable("MINIO_ENDPOINT")
	cfg.MINIO_ACCESS_KEY = loadEnvVariable("MINIO_ACCESS_KEY")
	cfg.MINIO_SECRET_KEY = loadEnvVariable("MINIO_SECRET_KEY")
	cfg.MINIO_SECURE = loadEnvVariableBool("MINIO_SECURE")
	cfg.MINIO_BUCKET = loadEnvVariable("MINIO_BUCKET")
	cfg.MINIO_WORKERS = loadEnvVariableInt("MINIO_WORKERS")
	return &cfg
}

func loadEnvVariableBool(key string) bool {
	value, exists := os.LookupEnv(key)
	parsedValue, err := strconv.ParseBool(value)

	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	if err != nil {
		log.Fatalf("Invalid value for %s: %v", key, err)
	}
	return parsedValue
}

func loadEnvVariableInt(key string) int {
	value, exists := os.LookupEnv(key)
	parsedValue, err := strconv.ParseInt(value, 10, 32)

	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	if err != nil {
		log.Fatalf("Invalid value for %s: %v", key, err)
	}
	return int(parsedValue)
}

func loadEnvVariable(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
