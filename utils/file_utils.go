package utils

import (
	"os"
)

func FileExists(filename string) bool {
	file, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !file.IsDir()
}

func FileNotEmpty(filename string) bool {
	file, err := os.Stat(filename)
	if err != nil {
		return false
	}

	size := file.Size()
	if size > 0 {
		return true
	}
	return false
}