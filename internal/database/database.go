package database

import (
	"context"
	"fmt"
	"log"

	"film-downloader/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(cfg *config.Config) *pgxpool.Pool {
	connectionString := buildConnectionString(cfg)
	dbCfg, err := pgxpool.ParseConfig(connectionString)
	dbCfg.MaxConns = cfg.DB_MAX_CONNS

	if err != nil {
		log.Fatalf("Unable to parse database config💊: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)

	if err != nil {
		log.Fatalf("failed to create connection poolpool🏊: %v\n", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Could not ping postgres🫙 database: %v", err))
	}

	log.Println("Database 🥳 connection pool initialized successfully ✅")
	return pool
}

func buildConnectionString(cfg *config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DB_USER, cfg.DB_PASSWORD,
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME,
	)
}
