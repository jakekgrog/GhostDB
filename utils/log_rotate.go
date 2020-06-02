package utils

import (
	"fmt"
	"io"
	"os"
	"os/user"
)

func LogMustRotate(logfile string, maxSize int64) (bool, error) {
	fi, err := os.Stat(logfile)
	if err != nil {
		return false, err
	}
	// get the size
	size := fi.Size()
	if size >= maxSize {
		return true, nil
	}
	return false, nil
}

func Rotate(filename string, tmpFilename string) (int64, error) {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	src := configPath + filename
	tmp := configPath + tmpFilename

	// Check if tmp file exists
	exists, err := tmpFileExists(tmp)
	if err != nil {
		return 0, fmt.Errorf("Error when checking for temp log existence: %s", err.Error())
	}

	// If it exists, clear it
	if exists {
		_, err := cleanFile(tmp)
		if err != nil {
			return 0, fmt.Errorf("failed to clean temp log")
		}
	}

	// Open the source file (main log)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, fmt.Errorf("failed to stat log file")
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	// Open the tmp log (or create if it doesn't exist)
	dst, err := os.OpenFile(tmp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s temporary log file", tmp)
	}
	defer dst.Close()

	// Copy the contents of main log to tmp log
	nBytes, err := io.Copy(dst, source)

	if err != nil {
		return 0, fmt.Errorf("failed to copy log to temp log")
	}

	// clear the main log
	_, err = cleanFile(src)

	if err != nil {
		return 0, fmt.Errorf("failed to clean snitch log file")
	}

	return nBytes, err
}

func tmpFileExists(tmpFilename string) (bool, error) {
	if _, err := os.Stat(tmpFilename); os.IsNotExist(err) {
		dst, err := os.OpenFile(tmpFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return false, fmt.Errorf("failed to open %s temporary log file", tmpFilename)
		}
		defer dst.Close()
	}
	return true, nil
}

func cleanFile(filePath string) (bool, error) {
	err := os.Remove(filePath)

	if err != nil {
		return false, err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	return true, err
}