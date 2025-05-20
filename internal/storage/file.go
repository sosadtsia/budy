package storage

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

// Storage defines the interface for storing and retrieving data
type Storage interface {
	Save(key string, data interface{}) error
	Load(key string, data interface{}) error
}

// FileStorage implements the Storage interface using files
type FileStorage struct {
	dataDir string
}

// NewFileStorage creates a new file storage
func NewFileStorage() (*FileStorage, error) {
	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Create a directory for our app if it doesn't exist
	dataDir := filepath.Join(usr.HomeDir, ".budy")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &FileStorage{
		dataDir: dataDir,
	}, nil
}

// GetDataDir returns the data directory path
func (s *FileStorage) GetDataDir() string {
	return s.dataDir
}

// Save stores data under the given key
func (s *FileStorage) Save(key string, data interface{}) error {
	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write to file
	filePath := filepath.Join(s.dataDir, key+".json")
	return os.WriteFile(filePath, jsonData, 0644)
}

// Load retrieves data stored under the given key
func (s *FileStorage) Load(key string, data interface{}) error {
	filePath := filepath.Join(s.dataDir, key+".json")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	// Read file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal JSON to data
	return json.Unmarshal(jsonData, data)
}
