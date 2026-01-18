package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"upnext/internal/model"
)

// JSONStore implements Store using a JSON file
type JSONStore struct {
	path string
}

// NewJSONStore creates a new JSON file store at the XDG-compliant path
func NewJSONStore() (*JSONStore, error) {
	path, err := getDataPath()
	if err != nil {
		return nil, err
	}
	return &JSONStore{path: path}, nil
}

// getDataPath returns the XDG-compliant path for the data file
func getDataPath() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "darwin":
		// macOS: Use ~/.local/share for consistency with Linux
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, ".local", "share")
	case "windows":
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(home, "AppData", "Local")
		}
	default:
		// Linux and others: Use XDG_DATA_HOME or fallback
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(home, ".local", "share")
		}
	}

	return filepath.Join(baseDir, "upnext", "todos.json"), nil
}

// Load reads the data from the JSON file
func (s *JSONStore) Load() (*model.Data, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return model.NewData(), nil
		}
		return nil, err
	}

	var result model.Data
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Save writes the data to the JSON file atomically
func (s *JSONStore) Save(data *model.Data) error {
	// Ensure directory exists
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first for atomic operation
	tempPath := s.path + ".tmp"
	if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
		return err
	}

	// Rename temp file to actual file (atomic on most systems)
	return os.Rename(tempPath, s.path)
}
