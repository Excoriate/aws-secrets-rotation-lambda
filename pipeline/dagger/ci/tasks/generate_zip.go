package tasks

import (
	"context"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/lambda"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
)

func GenerateZip() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ux.ShowSubTitle("lambda:", "Package")

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
	cfgValue, err := cfg.GetFromViper("lambda-src")
	lambdaSRCDir := cfgValue.Value.(string)

	// Validating lambda source directory.
	compiler := lambda.NewCompiler(client, ctx)
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

	defer client.Close()

	msg.ShowSuccess("", fmt.Sprintf("The zip file has been created: %s", zipPath))

	return nil
}
