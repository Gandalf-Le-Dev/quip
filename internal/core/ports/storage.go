package ports

import (
    "context"
    "io"
)

type Storage interface {
    Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
}