package ports

import (
	"context"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
)

type FileRepository interface {
    Store(ctx context.Context, file *domain.File) error
    FindByID(ctx context.Context, id string) (*domain.File, error)
    IncrementDownloads(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context) error
}

type PasteRepository interface {
    Store(ctx context.Context, paste *domain.Paste) error
    FindByID(ctx context.Context, id string) (*domain.Paste, error)
    IncrementViews(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context) error
}