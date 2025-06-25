package postgres

import (
	"context"
	"database/sql"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/ports"
	_ "github.com/lib/pq"
)

// FileRepository implementation
type Repository struct {
	db      *sql.DB
	queries *Queries
}

var _ ports.FileRepository = (*Repository)(nil)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		queries: New(db),
	}
}

func (r *Repository) Store(ctx context.Context, file *domain.File) error {
	_, err := r.queries.CreateFile(ctx, CreateFileParams{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		ContentType:  file.ContentType,
		StorageKey:   file.StorageKey,
		Downloads:    int32(file.Downloads),
		MaxDownloads: int32(file.MaxDownloads),
		CreatedAt:    file.CreatedAt,
		ExpiresAt:    file.ExpiresAt,
	})
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.File, error) {
	row, err := r.queries.GetFileByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &domain.File{
		ID:           row.ID,
		OriginalName: row.OriginalName,
		Size:         row.Size,
		ContentType:  row.ContentType,
		StorageKey:   row.StorageKey,
		Downloads:    int(row.Downloads),
		MaxDownloads: int(row.MaxDownloads),
		CreatedAt:    row.CreatedAt,
		ExpiresAt:    row.ExpiresAt,
	}, nil
}

func (r *Repository) IncrementDownloads(ctx context.Context, id string) error {
	return r.queries.IncrementFileDownloads(ctx, id)
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredFiles(ctx)
}

// PasteRepository implementation
type PasteRepository struct {
	*Repository
}

var _ ports.PasteRepository = (*PasteRepository)(nil)

func NewPasteRepository(db *sql.DB) ports.PasteRepository {
	return &PasteRepository{
		Repository: NewRepository(db),
	}
}

func (r *PasteRepository) Store(ctx context.Context, paste *domain.Paste) error {
	_, err := r.queries.CreatePaste(ctx, CreatePasteParams{
		ID:        paste.ID,
		Content:   paste.Content,
		Language:  paste.Language,
		Title:     sql.NullString{String: paste.Title, Valid: paste.Title != ""},
		Views:     int32(paste.Views),
		MaxViews:  int32(paste.MaxViews),
		CreatedAt: paste.CreatedAt,
		ExpiresAt: paste.ExpiresAt,
	})
	return err
}

func (r *PasteRepository) FindByID(ctx context.Context, id string) (*domain.Paste, error) {
	row, err := r.queries.GetPasteByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &domain.Paste{
		ID:        row.ID,
		Content:   row.Content,
		Language:  row.Language,
		Title:     row.Title.String,
		Views:     int(row.Views),
		MaxViews:  int(row.MaxViews),
		CreatedAt: row.CreatedAt,
		ExpiresAt: row.ExpiresAt,
	}, nil
}

func (r *PasteRepository) IncrementViews(ctx context.Context, id string) error {
	return r.queries.IncrementPasteViews(ctx, id)
}

func (r *PasteRepository) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredPastes(ctx)
}
