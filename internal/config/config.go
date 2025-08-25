package config

import (
	"film-downloader/internal/utils"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	accessToken  string
	tokenMutex   sync.RWMutex
	Cookie       string
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

	cfg.accessToken = utils.LoadString("ACCESS_TOKEN")
	cfg.Cookie = utils.LoadString("COOKIE")
	cfg.DB_HOST = utils.LoadString("DB_HOST")
	cfg.DB_PORT = utils.LoadString("DB_PORT")
	cfg.DB_USER = utils.LoadString("DB_USER")
	cfg.DB_PASSWORD = utils.LoadString("DB_PASSWORD")
	cfg.DB_NAME = utils.LoadString("DB_NAME")
	cfg.DB_MAX_CONNS = utils.LoadInt("DB_MAX_CONNS")
	cfg.SecureKey = "w3r1Sec4re_Token_"

	cfg.MINIO_ENDPOINT = utils.LoadString("MINIO_ENDPOINT")
	cfg.MINIO_ACCESS_KEY = utils.LoadString("MINIO_ACCESS_KEY")
	cfg.MINIO_SECRET_KEY = utils.LoadString("MINIO_SECRET_KEY")
	cfg.MINIO_SECURE = utils.LoadBool("MINIO_SECURE")
	cfg.MINIO_BUCKET = utils.LoadString("MINIO_BUCKET")
	cfg.MINIO_WORKERS = utils.LoadInt("MINIO_WORKERS")
	return &cfg
}

// GetAccessToken returns the access token in a thread-safe manner
func (c *Config) GetAccessToken() string {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()
	return c.accessToken
}

// SetAccessToken sets the access token in a thread-safe manner
func (c *Config) SetAccessToken(token string) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.accessToken = token
}
