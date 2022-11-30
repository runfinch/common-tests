package ffs

import (
	"os"
)

// CheckIfFileExists checks if a file exists in the filesystem.
func CheckIfFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}
