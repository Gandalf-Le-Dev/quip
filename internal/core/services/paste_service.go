package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/ports"
)

type PasteService struct {
	repo ports.PasteRepository
	log  *slog.Logger
}

func NewPasteService(repo ports.PasteRepository, log *slog.Logger) *PasteService {
	return &PasteService{
		repo: repo,
		log:  log,
	}
}

func (s *PasteService) Create(ctx context.Context, content, language, title string, ttl time.Duration) (*domain.Paste, error) {
	if content == "" {
		s.log.Warn("Attempt to create paste with empty content")
		return nil, domain.ErrInvalidInput
	}

	paste := domain.NewPaste(content, language, title, ttl)
	logger := s.log.With("paste_id", paste.ID)

	if err := s.repo.Store(ctx, paste); err != nil {
		logger.Error("Failed to store paste", "error", err)
		return nil, err
	}

	logger.Info("Paste created successfully", "language", paste.Language, "title", paste.Title)
	return paste, nil
}

func (s *PasteService) Get(ctx context.Context, id string) (*domain.Paste, error) {
	logger := s.log.With("paste_id", id)
	paste, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("Paste not found", "error", err)
		return nil, err
	}

	if !paste.CanView() {
		if paste.IsExpired() {
			logger.Warn("Attempt to view expired paste")
			return nil, domain.ErrExpired
		}
		logger.Warn("Attempt to view paste over limit")
		return nil, domain.ErrLimitExceeded
	}

	// Increment view counter
	if err := s.repo.IncrementViews(ctx, id); err != nil {
		logger.Error("Failed to increment paste views", "error", err)
		return nil, err
	}

	logger.Info("Paste viewed successfully")
	return paste, nil
}

func (s *PasteService) GetRaw(ctx context.Context, id string) (string, error) {
	paste, err := s.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return paste.Content, nil
}

func (s *PasteService) Delete(ctx context.Context, id string) error {
	// Implementation depends on your repository
	s.log.Warn("Delete paste not implemented", "paste_id", id)
	return nil
}

func (s *PasteService) CleanupExpired(ctx context.Context) error {
	s.log.Debug("Cleaning up expired pastes")
	err := s.repo.DeleteExpired(ctx)
	if err != nil {
		s.log.Error("Failed to cleanup expired pastes", "error", err)
	}
	return err
}
