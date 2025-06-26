package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/adapters/api"
	"github.com/Gandalf-Le-Dev/quip/internal/adapters/repository/postgres"
	"github.com/Gandalf-Le-Dev/quip/internal/adapters/repository/storage/minio"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
	"github.com/Gandalf-Le-Dev/quip/internal/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	log := logger.NewDevelopmentLogger()

	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fileshare:secretpassword@localhost:5432/fileshare?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Error("Failed to ping the database", "error", err)
		os.Exit(1)
	}
	log.Info("Successfully connected to the database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("Database migrations completed successfully")

	// Initialize repositories
	fileRepo := postgres.NewRepository(db)
	pasteRepo := postgres.NewPasteRepository(db)

	// Initialize storage
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "uploads"
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}

	storage, err := minio.NewMinioStorage(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		minioBucket,
		os.Getenv("MINIO_USE_SSL") == "true",
		log,
	)
	if err != nil {
		log.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}
	log.Info("Object storage initialized successfully")

	// Initialize services
	fileService := services.NewFileService(fileRepo, storage, log)
	pasteService := services.NewPasteService(pasteRepo, log)

	// Start cleanup goroutine
	go startCleanupTask(log, fileService, pasteService)

	// Initialize HTTP handlers
	handlers := api.NewHandlers(fileService, pasteService, log)
	router := api.NewRouter(handlers)

	// Start server
	log.Info("Server starting", "address", ":8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func startCleanupTask(log *slog.Logger, fileService *services.FileService, pasteService *services.PasteService) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Info("Running cleanup task for expired content")
		ctx := context.Background()

		if err := fileService.CleanupExpired(ctx); err != nil {
			log.Error("Error cleaning up expired files", "error", err)
		}

		if err := pasteService.CleanupExpired(ctx); err != nil {
			log.Error("Error cleaning up expired pastes", "error", err)
		}
		log.Info("Cleanup task finished")
	}
}
func runMigrations(db *sql.DB) error {
	// In a real application, use a migration tool like golang-migrate
	// For now, just execute the schema
	schema := `
		CREATE TABLE IF NOT EXISTS files (
		id VARCHAR(11) PRIMARY KEY,
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
		id VARCHAR(11) PRIMARY KEY,
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
