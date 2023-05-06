package common

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
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

func FileExist(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
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

func CreateZipFile(sourceFile, targetFile, targetDir string) (*os.File, error) {
	if sourceFile == "" {
		return nil, errors.New("source file cannot be an empty string")
	}

	if targetFile == "" {
		return nil, errors.New("target file cannot be an empty string")
	}

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file %s does not exist", sourceFile)
	}

	if targetDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("unable to get current directory: %s", err.Error())
		}
		targetDir = currentDir
	}

	zipFile, err := os.Create(filepath.Join(targetDir, targetFile))
	if err != nil {
		return nil, fmt.Errorf("error creating zip file %s: %s", targetFile, err.Error())
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	inputFile, err := os.Open(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("error opening source file %s: %s", sourceFile, err.Error())
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info for %s: %s", sourceFile, err.Error())
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return nil, fmt.Errorf("error creating zip header for %s: %s", sourceFile, err.Error())
	}
	header.Name = sourceFile

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return nil, fmt.Errorf("error creating zip writer for %s: %s", sourceFile, err.Error())
	}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return nil, fmt.Errorf("error writing to zip file %s: %s", sourceFile, err.Error())
	}

	return zipFile, nil
}
