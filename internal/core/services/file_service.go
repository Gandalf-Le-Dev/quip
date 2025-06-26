package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/ports"
)

type FileService struct {
	repo    ports.FileRepository
	storage ports.Storage
	log     *slog.Logger
}

func NewFileService(repo ports.FileRepository, storage ports.Storage, log *slog.Logger) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
		log:     log,
	}
}

func (s *FileService) Upload(ctx context.Context, reader io.Reader, filename string, size int64, contentType string, ttl time.Duration) (*domain.File, error) {
	file := domain.NewFile(filename, size, contentType, ttl)
	logger := s.log.With("file_id", file.ID, "storage_key", file.StorageKey)

	// Upload to storage
	if err := s.storage.Upload(ctx, file.StorageKey, reader, size, contentType); err != nil {
		logger.Error("Failed to upload file to storage", "error", err)
		return nil, err
	}
	logger.Debug("File uploaded to storage")

	// Save metadata to repository
	if err := s.repo.Store(ctx, file); err != nil {
		// Cleanup storage on failure
		logger.Error("Failed to store file metadata, cleaning up storage", "error", err)
		if delErr := s.storage.Delete(ctx, file.StorageKey); delErr != nil {
			logger.Error("Failed to cleanup storage after metadata store failure", "delete_error", delErr)
		}
		return nil, err
	}
	logger.Info("File uploaded successfully", "file_id", file.ID, "storage_key", file.StorageKey)

	return file, nil
}

func (s *FileService) Download(ctx context.Context, id string) (io.ReadCloser, *domain.File, error) {
	logger := s.log.With("file_id", id)
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("File not found in repository", "error", err)
		return nil, nil, err
	}

	logger.Debug("File found in repository",
		"file_id", file.ID,
		"original_name", file.OriginalName,
		"size", file.Size,
		"content_type", file.ContentType,
		"storage_key", file.StorageKey,
		"downloads", file.Downloads,
		"max_downloads", file.MaxDownloads,
		"created_at", file.CreatedAt,
		"expires_at", file.ExpiresAt,
		"is_expired", file.IsExpired(),
		"can_download", file.CanDownload(),
	)

	logger.Debug("Attempting to download from storage", "storage_key", file.StorageKey)

	if !file.CanDownload() {
		if file.IsExpired() {
			logger.Warn("Attempt to download expired file")
			return nil, nil, domain.ErrExpired
		}
		logger.Warn("Attempt to download file over limit")
		return nil, nil, domain.ErrLimitExceeded
	}

	reader, err := s.storage.Download(ctx, file.StorageKey)
	if err != nil {
		logger.Error("Failed to download file from storage", "storage_key", file.StorageKey, "error", err)
		return nil, nil, err
	}
	logger.Debug("File downloaded from storage")

	// Increment download counter
	if err := s.repo.IncrementDownloads(ctx, id); err != nil {
		reader.Close()
		logger.Error("Failed to increment file download count", "error", err)
		return nil, nil, err
	}
	logger.Info("File downloaded successfully")

    file.OriginalName = appendTimestamp(file.OriginalName)

	return reader, file, nil
}

func (s *FileService) GetInfo(ctx context.Context, id string) (*domain.File, error) {
	s.log.Debug("Fetching file info", "file_id", id)
	return s.repo.FindByID(ctx, id)
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	logger := s.log.With("file_id", id)
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Error("Failed to find file for deletion", "error", err)
		return err
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, file.StorageKey); err != nil {
		logger.Error("Failed to delete file from storage", "storage_key", file.StorageKey, "error", err)
		return err
	}
	logger.Debug("File deleted from storage")

	// Delete from repository
	// Implementation depends on your repository
	logger.Info("File deleted successfully")
	return nil
}

func (s *FileService) CleanupExpired(ctx context.Context) error {
	s.log.Debug("Cleaning up expired files")
	err := s.repo.DeleteExpired(ctx)
	if err != nil {
		s.log.Error("Failed to cleanup expired files", "error", err)
	}
	return err
}
func appendTimestamp(filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	timestamp := now.Format("20060102_150405")
	return fmt.Sprintf("%s_%s%s", name, timestamp, ext)
}
