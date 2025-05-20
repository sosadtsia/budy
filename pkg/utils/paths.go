package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands ~ to user's home directory
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}
	return path, nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// EnsureDirExists creates a directory if it doesn't exist
func EnsureDirExists(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// GetDataDir returns the application data directory
func GetDataDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(usr.HomeDir, ".budy")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}
