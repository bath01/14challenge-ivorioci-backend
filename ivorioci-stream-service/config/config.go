package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	GoEnv                string
	DatabaseURL          string
	JWTAccessSecret      string
	VideoStoragePath     string
	ThumbnailStoragePath string
	PublicBaseURL        string // e.g. http://localhost:3000/api — used to build thumbnail URLs
}

func Load() *Config {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, reading from environment")
		}
	}

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		GoEnv:                getEnv("GO_ENV", "development"),
		DatabaseURL:          mustGetEnv("DATABASE_URL"),
		JWTAccessSecret:      mustGetEnv("JWT_ACCESS_SECRET"),
		VideoStoragePath:     getEnv("VIDEO_STORAGE_PATH", "./storage/videos"),
		ThumbnailStoragePath: getEnv("THUMBNAIL_STORAGE_PATH", "./storage/thumbnails"),
		PublicBaseURL:        getEnv("PUBLIC_BASE_URL", "http://localhost:3000/api"),
	}
}

func InitDB(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("[DB] Connected to PostgreSQL")
	return pool, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[CONFIG] Required environment variable %q is not set", key)
	}
	return val
}
