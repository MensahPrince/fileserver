package storage

import (
	"context"
	"errors"
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

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // File or directory exists
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Doesn't exist
	}

	// Some other error (e.g. permission denied)
	return false, err
}

func (d *DiskStorage) Delete(ctx context.Context, id string) error {
	path, err := d.resolvePath(id)
	if err != nil {
		return err
	}

	//Delete data
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("Failed to remove data")
	}

	return nil
}
