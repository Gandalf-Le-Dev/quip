package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/adapters/repository/postgres"
	"github.com/Gandalf-Le-Dev/quip/internal/adapters/repository/storage/minio"
	"github.com/Gandalf-Le-Dev/quip/internal/adapters/web"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/fileshare?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	fileRepo := postgres.NewRepository(db)
	pasteRepo := postgres.NewPasteRepository(db)

	// Initialize storage
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	storage, err := minio.NewMinioStorage(
		minioEndpoint,
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_BUCKET"),
		os.Getenv("MINIO_USE_SSL") == "true",
	)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// Initialize services
	fileService := services.NewFileService(fileRepo, storage)
	pasteService := services.NewPasteService(pasteRepo)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			ctx := context.Background()
			err := fileService.CleanupExpired(ctx)
			if err != nil {
				log.Println("Error cleaning up expired files:", err)
				continue
			}
			err = pasteService.CleanupExpired(ctx)
			if err != nil {
				log.Println("Error cleaning up expired pastes:", err)
				continue
			}
		}
	}()

	// Initialize HTTP handlers
	handlers := web.NewHandlers(fileService, pasteService)
	router := web.NewRouter(handlers)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func runMigrations(db *sql.DB) error {
	// In a real application, use a migration tool like golang-migrate
	// For now, just execute the schema
	schema := `
    CREATE TABLE IF NOT EXISTS files (
        id VARCHAR(10) PRIMARY KEY,
        original_name VARCHAR(255) NOT NULL,
        size BIGINT NOT NULL,
        content_type VARCHAR(100) NOT NULL,
        storage_key VARCHAR(100) NOT NULL UNIQUE,
        downloads INT NOT NULL DEFAULT 0,
        max_downloads INT NOT NULL DEFAULT -1,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        expires_at TIMESTAMP NOT NULL
    );
    
    CREATE TABLE IF NOT EXISTS pastes (
        id VARCHAR(10) PRIMARY KEY,
        content TEXT NOT NULL,
        language VARCHAR(50) NOT NULL,
        title VARCHAR(255),
        views INT NOT NULL DEFAULT 0,
        max_views INT NOT NULL DEFAULT -1,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        expires_at TIMESTAMP NOT NULL
    );
    
    CREATE INDEX IF NOT EXISTS idx_files_expires_at ON files(expires_at);
    CREATE INDEX IF NOT EXISTS idx_pastes_expires_at ON pastes(expires_at);
    `

	_, err := db.Exec(schema)
	return err
}
