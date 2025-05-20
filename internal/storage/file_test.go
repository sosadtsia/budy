package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "budy-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temp directory: %v", err)
		}
	}()

	// Create a test storage instance with the temp directory
	storage := &FileStorage{
		dataDir: tempDir,
	}

	// Test GetDataDir method
	t.Run("GetDataDir", func(t *testing.T) {
		dir := storage.GetDataDir()
		if dir != tempDir {
			t.Errorf("GetDataDir returned wrong directory. Got: %s, Expected: %s", dir, tempDir)
		}
	})

	// Test Save and Load methods
	t.Run("SaveAndLoad", func(t *testing.T) {
		// Test data
		type TestData struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		testKey := "test-key"
		testData := TestData{
			Name:  "Test",
			Value: 42,
		}

		// Save data
		err := storage.Save(testKey, testData)
		if err != nil {
			t.Fatalf("Failed to save data: %v", err)
		}

		// Verify file exists
		filePath := filepath.Join(tempDir, testKey+".json")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("File was not created at %s", filePath)
		}

		// Load data
		var loadedData TestData
		err = storage.Load(testKey, &loadedData)
		if err != nil {
			t.Fatalf("Failed to load data: %v", err)
		}

		// Verify data integrity
		if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
			t.Errorf("Loaded data doesn't match original. Got: %+v, Expected: %+v", loadedData, testData)
		}
	})

	// Test loading non-existent data
	t.Run("LoadNonExistent", func(t *testing.T) {
		var data map[string]interface{}
		err := storage.Load("non-existent", &data)
		if err == nil {
			t.Error("Expected an error when loading non-existent key, but got nil")
		}
	})
}
