package utils

import (
	"fmt"
	"path/filepath"
)

func Traversal(path string) bool {

	fmt.Printf("path is: %s", path)
	//clean the path
	clean_path := filepath.Clean(path)
	fmt.Printf("Cleaned Path: %s ", clean_path)
	//Check if path is absolute

	return filepath.IsLocal(clean_path)
}
