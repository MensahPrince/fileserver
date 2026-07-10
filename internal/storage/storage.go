package storage

import (
	"context"
	"io"
)

type Storage interface {
    Save(ctx context.Context, id string, r io.Reader) (int64, error)
    Open(ctx context.Context, id string) (io.ReadCloser, error)
    Delete(ctx context.Context, id string) error
}
