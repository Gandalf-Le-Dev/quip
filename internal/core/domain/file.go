package domain

import (
	"fmt"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/pkg/utils/nanoid"
)

type File struct {
	ID           string
	OriginalName string
	Size         int64
	ContentType  string
	StorageKey   string
	Downloads    int
	MaxDownloads int
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func NewFile(originalName string, size int64, contentType string, ttl time.Duration) *File {
	return &File{
		ID:           generateID(),
		OriginalName: originalName,
		Size:         size,
		ContentType:  contentType,
		StorageKey:   generateStorageKey(),
		Downloads:    0,
		MaxDownloads: -1, // unlimited
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(ttl),
	}
}

func (f *File) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

func (f *File) CanDownload() bool {
	if f.IsExpired() {
		return false
	}
	if f.MaxDownloads > 0 && f.Downloads >= f.MaxDownloads {
		return false
	}
	return true
}

func (f *File) IncrementDownloads() {
	f.Downloads++
}

func generateID() string {
	// Generate short, URL-safe ID with nanoid
	return nanoid.Must(11)
}

func generateStorageKey() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), generateID())
}
