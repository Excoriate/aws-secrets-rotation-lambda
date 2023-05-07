package tasks

import (
	"context"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/lambda"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
	"path/filepath"
)

func PackageZip() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ux.ShowSubTitle("lambda:", "PackageZIP")

	// Getting working directories.
	dirs, err := config.GetDirConfig()
	if err != nil {
		msg.ShowError("", "Failed to get working directories", err)
		return err
	}

	msg.ShowInfo("", fmt.Sprintf("Current root directory: %s", dirs.CurrentDir))

	// Booting dagger!
	ctx := context.Background()
	client, err := daggerio.NewClient(ctx)
	if err != nil {
		msg.ShowError("", "Failed to boot dagger", err)
		return err
	}

	// Fetching configuration from Viper.
	cfg := config.Cfg{}

	// Check if there's an existing binary that's passed. If so, use it.
	existingBinaryCfg, err := cfg.GetFromViper("existing-binary")
	var lambdaSRCDir string
	var compiledBinaryPath string
	var existingBinary bool
	if err == nil {
		compiledBinaryPath = existingBinaryCfg.Value.(string)
		compiledBinaryPath = filepath.Join(dirs.CurrentDir, compiledBinaryPath)

		msg.ShowWarning("", fmt.Sprintf("Using existing binary: %s", compiledBinaryPath))

		if err := common.FileExist(compiledBinaryPath); err != nil {
			msg.ShowError("", "The existing binary does not exist", err)
			return err
		}

		lambdaSRCDir = compiledBinaryPath
		existingBinary = true

	} else {
		cfgValue, err := cfg.GetFromViper("lambda-src")
		if err != nil {
			msg.ShowError("", "Failed to get lambda source directory from Viper", err)
			return err
		}
		lambdaSRCDir = cfgValue.Value.(string)
		existingBinary = false
	}

	compiler := lambda.NewCompiler(client, ctx)

	if !existingBinary {
		// Validating lambda source directory.
		srcPath, err := lambda.IsLambdaSRCDirValid(lambdaSRCDir, dirs)

		if err != nil {
			msg.ShowError("", "Failed to validate lambda source directory", err)
			return err
		}

		// Compiling
		_, binaryExportedPath, err := compiler.Compile(srcPath, "")
		if err != nil {
			msg.ShowError("", "The binary could not be compiled", err)
			return err
		}

		// Output paths
		srcFileToZip := fmt.Sprintf("%s/%s", config.GetBinaryExportPath(), binaryExportedPath)
		targetFile := config.PackageZipName
		targetDir := config.GetZipExportPath()

		// Creating the zip file from the binary.
		_, zipPath, err := compiler.Zip(srcFileToZip, targetFile, targetDir)
		if err != nil {
			msg.ShowError("", "Failed to create zip file", err)
			return err
		}
		msg.ShowSuccess("", fmt.Sprintf("The zip file has been created: %s", zipPath))
	} else {
		srcFileToZip := compiledBinaryPath
		targetFile := config.PackageZipName
		targetDir := config.GetZipExportPath()

		// Creating the zip file from the binary.
		_, zipPath, err := compiler.Zip(srcFileToZip, targetFile, targetDir)
		if err != nil {
			msg.ShowError("", "Failed to create zip file", err)
			return err
		}
		msg.ShowSuccess("", fmt.Sprintf("The zip file has been created: %s", zipPath))
	}

	defer client.Close()

	return nil
}
