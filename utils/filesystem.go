package utils

import (
	"fmt"
	"os"
)

func EnsureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return fmt.Errorf("could not create directory %v. %v", dir, err)
		}
	}
	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func WriteFile(filePath string, contents []byte, mode os.FileMode) error {
	return os.WriteFile(filePath, contents, mode)
}

func MustWriteFile(filePath string, contents []byte, mode os.FileMode) error {
	err := WriteFile(filePath, contents, mode)
	if err != nil {
		return fmt.Errorf("MustWriteFile failed: %v", err)
	}
	return nil
}
