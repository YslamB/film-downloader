package main

import (
	"database/sql"
	"film-downloader/internal/config"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func connectToClickHouse(cfg *config.Config) (*sql.DB, error) {
	dsn := "clickhouse://default:S3ku4!wor6@95.85.126.219:9000/bb_analytics?dial_timeout=2s&compress=true"

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	log.Println("✅ Successfully connected to ClickHouse")
	return db, nil
}

func main() {
	cfg := config.Init()
	db, err := connectToClickHouse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// sudo ufw allow proto tcp from 93.171.220.16 to any port 5469
// sudo ufw allow proto tcp from 93.171.220.16 to any port 9000
// sudo ufw allow proto tcp from 192.168.10.0/24 to any port 9000
// sudo ufw allow proto tcp from 192.168.10.0/24 to any port 5469
