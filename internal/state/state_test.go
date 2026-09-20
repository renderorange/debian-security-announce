package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LastDSANumber != 0 {
		t.Errorf("got last DSA number %d, want 0", s.LastDSANumber)
	}
}

func TestLoad_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	data := []byte(`{"last_dsa_number": 6500, "last_checked": "2026-09-20T00:00:00Z"}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LastDSANumber != 6500 {
		t.Errorf("got last DSA number %d, want 6500", s.LastDSANumber)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	data := []byte(`{invalid json`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s := &State{
		LastDSANumber: 6508,
		LastChecked:   "2026-09-20T03:24:14Z",
	}

	if err := Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("failed to load saved state: %v", err)
	}
	if loaded.LastDSANumber != 6508 {
		t.Errorf("got last DSA number %d, want 6508", loaded.LastDSANumber)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "state.json")

	s := &State{LastDSANumber: 100, LastChecked: "2026-09-20T00:00:00Z"}

	if err := Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}
