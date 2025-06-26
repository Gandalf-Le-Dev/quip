package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/ports"
)

type FileService struct {
	repo    ports.FileRepository
	storage ports.Storage
}

func NewFileService(repo ports.FileRepository, storage ports.Storage) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

func (s *FileService) Upload(ctx context.Context, reader io.Reader, filename string, size int64, contentType string, ttl time.Duration) (*domain.File, error) {
	file := domain.NewFile(filename, size, contentType, ttl)

	// Upload to storage
	if err := s.storage.Upload(ctx, file.StorageKey, reader, size, contentType); err != nil {
		return nil, err
	}

	// Save metadata to repository
	if err := s.repo.Store(ctx, file); err != nil {
		// Cleanup storage on failure
		_ = s.storage.Delete(ctx, file.StorageKey)
		return nil, err
	}

	return file, nil
}

func (s *FileService) Download(ctx context.Context, id string) (io.ReadCloser, *domain.File, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if !file.CanDownload() {
		if file.IsExpired() {
			return nil, nil, domain.ErrExpired
		}
		return nil, nil, domain.ErrLimitExceeded
	}

	reader, err := s.storage.Download(ctx, file.StorageKey)
	if err != nil {
		return nil, nil, err
	}

	// Increment download counter
	if err := s.repo.IncrementDownloads(ctx, id); err != nil {
		reader.Close()
		return nil, nil, err
	}

    file.OriginalName = appendTimestamp(file.OriginalName)

	return reader, file, nil
}

func (s *FileService) GetInfo(ctx context.Context, id string) (*domain.File, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, file.StorageKey); err != nil {
		return err
	}

	// Delete from repository
	// Implementation depends on your repository
	return nil
}

func (s *FileService) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}
func appendTimestamp(filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	timestamp := now.Format("20060102_150405")
	return fmt.Sprintf("%s_%s%s", name, timestamp, ext)
}
