package storage

import "testing"

func TestResolvePath(t *testing.T) {
	roots := []string{"./data", "/var/lib/fileserver"}
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"legit uuid", "7f3a9c2e-6b41-4d88-a5f7-9e2c1b6d4a90", false},
		{"parent traversal", "../../etc/passwd", true},
		{"single parent escape", "../main.go", true},
		{"absolute path", "/etc/passwd", true},
		{"empty id", "", true},
	}

	for _, r := range roots {
		d := NewDiskStorage(r)
		for _, tt := range tests {
			t.Run(r+"/"+tt.name, func(t *testing.T) {
				got, err := d.resolvePath(tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("resolvePath(%q) = %q, %v; wantErr %v", tt.id, got, err, tt.wantErr)
				}
			})
		}
	}
}
