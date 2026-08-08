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

var ErrNotFound = errors.New("cannot find file")

func (d *DiskStorage) resolvePath(id string) (string, error) {
	path := filepath.Join(d.root, id)
	cleanedPath := filepath.Clean(path)

	if !filepath.IsLocal(cleanedPath) {
		return "", fmt.Errorf("path escaped")
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
		return n, fmt.Errorf("failed to save: %w", err)
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

	err = os.Remove(path)

	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

func (d *DiskStorage) Open(ctx context.Context, id string) (io.ReadCloser, error) {
	path, err := d.resolvePath(id)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open: %w", err)
	}

	return file, nil
}
