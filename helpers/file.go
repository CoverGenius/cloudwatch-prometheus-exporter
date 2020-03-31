package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadFile returns the byte contents of a file located at path
//
// If an error is encountered while reading the file it is logged NOT returned
func ReadFile(path *string) *[]byte {
	absolutePath, err := filepath.Abs(*path)
	LogError(err)
	content, err := ioutil.ReadFile(absolutePath)
	LogError(err)
	return &content
}

// IsFileExists returns true if a file located a path exists
func IsFileExists(path *string) bool {
	_, err := os.Stat(*path)
	if err == nil {
		return true
	}
	return false
}
