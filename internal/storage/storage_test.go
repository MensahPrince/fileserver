package storage

import (
	"fmt"
	"path/filepath"
	"testing"
)

var root string = "./data"

func Traversal(path string) bool {

	fmt.Printf("path is: %s", path)
	//clean the path
	clean_path := filepath.Clean(path)
	fmt.Printf("Cleaned Path: %s ", clean_path)
	//Check if path is absolute

	return filepath.IsLocal(clean_path)
}

func TestTraversal(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantLocal bool
	}{
		{"legit uuid", "7f3a9c2e-6b41-4d88-a5f7-9e2c1b6d4a90", true},
		{"parent traversal", "../../etc/passwd", false},
		{"empty path", "", true},
		{"standard path", "C:/Users/Public/Documents/Project%20hiles/2026%20Reports/summary%20inal.pdf", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := filepath.Join(root, tt.id)
			got := Traversal(p)
			if got != tt.wantLocal {
				t.Errorf("Traversal(%q) = %v, want %v", p, got, tt.wantLocal)
			}
		})
	}
}
