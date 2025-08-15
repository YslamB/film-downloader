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
	DB_MAX_CONNS int32
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
	cfg.DB_MAX_CONNS = loadEnvVariableInt32("DB_MAX_CONNS")
	cfg.SecureKey = "w3r1Sec4re_Token_"
	return &cfg
}

func loadEnvVariableInt32(key string) int32 {
	value, exists := os.LookupEnv(key)
	parsedValue, err := strconv.ParseInt(value, 10, 32)

	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	if err != nil {
		log.Fatalf("Invalid value for %s: %v", key, err)
	}
	return int32(parsedValue)
}

func loadEnvVariable(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
