package config

import (
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/erroer"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const LambdaName = "secrets-manager-rotator-lambda"
const PackageZipName = "secrets-manager-rotator-lambda.zip"
const OutputBinaryDir = "output/lambda-bin"
const OutputZipDir = "output/lambda-zip"
const infraBaseDir = "infra/terraform"

var (
	CurrentDir, _      = os.Getwd()
	GitRelativeRepoDir = fmt.Sprintf("%s/../../", CurrentDir)
)

type DirConfig struct {
	CurrentDir             string
	GitRepoDir             string
	InfraDirBase           string
	GitRelativeRepoDir     string
	InfraDirModuleBucket   string
	InfraDirModuleFunction string
}

func GetBinaryExportPath() string {
	dirs, _ := GetDirConfig()
	// Replace trailing slash in the path
	gitRepoDirNormalised := strings.TrimSuffix(dirs.GitRepoDir, "/")
	return fmt.Sprintf("%s/%s", gitRepoDirNormalised, OutputBinaryDir)
}

func GetZipExportPath() string {
	dirs, _ := GetDirConfig()
	gitRepoDirNormalised := strings.TrimSuffix(dirs.GitRepoDir, "/")
	return fmt.Sprintf("%s/%s", gitRepoDirNormalised, OutputZipDir)
}

func GetInfraBaseDirPath() string {
	dirs, _ := GetDirConfig()
	gitRepoDirNormalised := strings.TrimSuffix(dirs.GitRepoDir, "/")
	return fmt.Sprintf("%s/%s", gitRepoDirNormalised, infraBaseDir)
}

func GetInfraWorkDirPath(module string) string {
	var path string

	if module == "bucket" {
		path = fmt.Sprintf("%s/%s", infraBaseDir, "lambda-deployment-bucket")
	}

	if module == "function" {
		path = fmt.Sprintf("%s/%s", infraBaseDir, "lambda-function")
	}

	if module == "secret" {
		path = fmt.Sprintf("%s/%s", infraBaseDir, "secrets-manager-secret")
	}

	if module == "data" {
		path = fmt.Sprintf("%s/%s", infraBaseDir, "lambda-data")
	}

	return path
}

func GetDirConfig() (DirConfig, error) {
	cfg := Cfg{}
	isDebugModeEnabled, _ := cfg.GetFromViperBool("debug")
	infraDirModuleBucket := fmt.Sprintf("%s/%s", infraBaseDir, "lambda-deployment-bucket")
	infraDirModuleFunction := fmt.Sprintf("%s/%s", infraBaseDir, "lambda-function")

	if isDebugModeEnabled {
		return DirConfig{
			CurrentDir:             CurrentDir,
			GitRepoDir:             GitRelativeRepoDir,
			InfraDirBase:           infraBaseDir,
			InfraDirModuleBucket:   infraDirModuleBucket,
			InfraDirModuleFunction: infraDirModuleFunction,
		}, nil
	}

	gitRootDir, err := GetGitRootDir()

	if err != nil {
		return DirConfig{}, erroer.NewPipelineConfigurationError(
			"failed to get the root of the Git repository", err)
	}
	return DirConfig{
		CurrentDir:             CurrentDir,
		GitRepoDir:             gitRootDir,
		InfraDirBase:           infraBaseDir,
		InfraDirModuleBucket:   infraDirModuleBucket,
		InfraDirModuleFunction: infraDirModuleFunction,
	}, nil
}

func GetGitRootDirRelative() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel", "../../../")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to find the root of the Git repository: %v", err)
	}

	gitRootPath := strings.TrimSpace(string(output))
	absolutePath, err := filepath.Abs(gitRootPath)

	if err != nil {
		return "", fmt.Errorf("failed to convert the Git root path to an absolute path: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current working directory: %v", err)
	}

	if currentDir == absolutePath {
		return absolutePath, nil
	}

	return "", fmt.Errorf("current directory is not the root of the Git repository")
}

func GetGitRootDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to find the root of the Git repository: %v", err)
	}

	gitRootPath := strings.TrimSpace(string(output))
	absolutePath, err := filepath.Abs(gitRootPath)

	if err != nil {
		return "", fmt.Errorf("failed to convert the Git root path to an absolute path: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current working directory: %v", err)
	}

	if currentDir == absolutePath {
		return absolutePath, nil
	}

	return "", fmt.Errorf("current directory is not the root of the Git repository")
}
