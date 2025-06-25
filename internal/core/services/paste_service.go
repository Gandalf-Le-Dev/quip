package services

import (
	"context"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/ports"
)

type PasteService struct {
	repo ports.PasteRepository
}

func NewPasteService(repo ports.PasteRepository) *PasteService {
	return &PasteService{
		repo: repo,
	}
}

func (s *PasteService) Create(ctx context.Context, content, language, title string, ttl time.Duration) (*domain.Paste, error) {
	if content == "" {
		return nil, domain.ErrInvalidInput
	}

	paste := domain.NewPaste(content, language, title, ttl)

	if err := s.repo.Store(ctx, paste); err != nil {
		return nil, err
	}

	return paste, nil
}

func (s *PasteService) Get(ctx context.Context, id string) (*domain.Paste, error) {
	paste, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !paste.CanView() {
		if paste.IsExpired() {
			return nil, domain.ErrExpired
		}
		return nil, domain.ErrLimitExceeded
	}

	// Increment view counter
	if err := s.repo.IncrementViews(ctx, id); err != nil {
		return nil, err
	}

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
	return nil
}

func (s *PasteService) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}
