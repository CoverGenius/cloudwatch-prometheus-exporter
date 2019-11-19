package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFile(path *string) *[]byte {
	absolutePath, err := filepath.Abs(*path)
	LogError(err)
	content, err := ioutil.ReadFile(absolutePath)
	LogError(err)
	return &content
}

func IsFileExists(path *string) bool {
	_, err := os.Stat(*path)
	if err == nil {
		return true
	} else {
		return false
	}
}
