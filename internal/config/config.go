package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessToken string
	SecureKey   string
}

func Init() Config {
	var cfg Config
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.AccessToken = os.Getenv("ACCESS_TOKEN")
	cfg.SecureKey = "w3r1Sec4re_Token_"
	return cfg
}
