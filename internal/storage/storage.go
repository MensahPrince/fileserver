package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(ctx context.Context, id string, r io.Reader) (int64, error)
	Open(ctx context.Context, id string) (io.ReadCloser, error)
	Delete(ctx context.Context, id string) error
}

type DiskStorage struct {
	root string
}

func NewDiskStorage(root string) *DiskStorage {
	return &DiskStorage{root: root}
}

func (d *DiskStorage) resolvePath(id string) (string, error) {
	path := filepath.Join(d.root, id)
	cleanedPath := filepath.Clean(path)
	 
	if !filepath.IsLocal(cleanedPath) {
		return "", fmt.Errorf("Path Escaped")
	}

	return cleanedPath, nil
}

func (d *DiskStorage) Save(ctx context.Context, id string, r io.Reader) (int64, error) {
	path, err := d.resolvePath(id)
	if err != nil {
		return 0, err
	}

	// os.Create makes a new file (or truncates an existing one with the same name) and gives you an *os.File, which itself implements io.Writer. Note the defer f.Close() — you want the file handle released once Save returns, regardless of success or failure.
	// Create file
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return n, err
	}

	return n, nil
}
