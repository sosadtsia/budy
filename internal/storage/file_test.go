package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestFileStorage tests the file storage implementation
func TestFileStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "budy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temporary directory: %v", err)
		}
	}()

	// Create a test storage with our temp directory
	storage := &FileStorage{
		dataDir: tempDir,
	}

	// Test GetDataDir
	if storage.GetDataDir() != tempDir {
		t.Errorf("Expected data dir to be %s, got %s", tempDir, storage.GetDataDir())
	}

	// Test data to save and load
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	testData := TestData{
		Name:  "test",
		Value: 42,
	}

	// Test Save
	key := "test-key"
	if err := storage.Save(key, testData); err != nil {
		t.Errorf("Error saving data: %v", err)
	}

	// Verify the file was created
	filePath := filepath.Join(tempDir, key+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", filePath)
	}

	// Verify file contents
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	var savedData TestData
	if err := json.Unmarshal(fileData, &savedData); err != nil {
		t.Errorf("Error unmarshaling file data: %v", err)
	}
	if savedData.Name != testData.Name || savedData.Value != testData.Value {
		t.Errorf("File data doesn't match expected: got %+v, want %+v", savedData, testData)
	}

	// Test Load
	var loadedData TestData
	if err := storage.Load(key, &loadedData); err != nil {
		t.Errorf("Error loading data: %v", err)
	}
	if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
		t.Errorf("Loaded data doesn't match expected: got %+v, want %+v", loadedData, testData)
	}

	// Test Load with non-existent key
	var emptyData TestData
	err = storage.Load("non-existent", &emptyData)
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("Expected not exists error, got %v", err)
	}
}
