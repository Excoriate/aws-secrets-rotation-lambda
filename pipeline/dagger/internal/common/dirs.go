package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func DirExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}

	return nil
}

func PathExist(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path %s does not exist", path)
		}
		return nil, fmt.Errorf("error checking path %s: %s", path, err.Error())
	}

	return info, nil
}

func PathIsADirectory(path string) error {
	info, err := PathExist(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}

func DirIsNotEmpty(dir string) error {
	if err := DirExist(dir); err != nil {
		return err
	}

	if err := PathIsADirectory(dir); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		return fmt.Errorf("failed to read content of directory %s", dir)
	}

	if len(entries) == 0 {
		return fmt.Errorf("directory %s is empty", dir)
	}

	return nil
}

func DirIsValid(dir string) error {
	if err := DirExist(dir); err != nil {
		return err
	}

	if err := PathIsADirectory(dir); err != nil {
		return err
	}

	return nil
}

func GetDirAbsolute(dir string) (string, error) {
	absolutePath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("error converting path %s to absolute path: %s", dir, err.Error())
	}

	return absolutePath, nil
}
