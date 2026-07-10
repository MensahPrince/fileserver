package storage

import (
	"context"
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

func (d *DiskStorage) Save(ctx context.Context, id string, r io.Reader) (int64, error) {
    path := filepath.Join(d.root, id)

    
    return 0, nil // Declared Zero for now, would return value after the function is done 
}