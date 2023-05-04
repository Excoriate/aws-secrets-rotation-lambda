package tasks

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
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

	// Fetching configuration from Viper.
	cfg := config.Cfg{}
	cfgValue, err := cfg.GetFromViper("lambda-src")
	var srcPath string

	// Validating the Lambda source code directory.
	if err != nil {
		msg.ShowWarning("", fmt.Sprintf("The lambda source code directory is not set. Using the current directory: %s", dirs.CurrentDir))
		srcPath = dirs.CurrentDir
	} else {
		srcPath = cfgValue.Value.(string)

		if err := common.DirIsValid(srcPath); err != nil {
			msg.ShowError("",
				fmt.Sprintf("The lambda source code directory is not a valid directory"+
					": %s. Current directory: %s", srcPath, dirs.CurrentDir), err)

			return err
		}

		if err := common.DirIsNotEmpty(srcPath); err != nil {
			msg.ShowError("", fmt.Sprintf("The lambda source code directory is empty: %s", srcPath), err)
			return err
		}

		absDirPath, _ := common.GetDirAbsolute(srcPath)
		srcPath = absDirPath

		msg.ShowInfo("", fmt.Sprintf("Using the lambda source code directory: %s", srcPath))
	}

	// Booting dagger!
	ctx := context.Background()
	client, err := daggerio.NewClient(ctx)
	if err != nil {
		return err
	}

	// Setting dagger 'workidr' and 'binary (output) directory'.
	workDir := client.Host().Directory(srcPath)

	// Output dir
	binaryOutDir := client.Directory()

	// Compiling the binary.
	platforms := []dagger.Platform{"linux/amd64"}
	for _, platform := range platforms {
		msg.ShowInfo("", fmt.Sprintf("Compiling the binary for platform: %s", platform))
		image := "golang:1.20"
		builder := client.Container(dagger.ContainerOpts{Platform: platform}).
			From(image).
			WithMountedDirectory("/src", workDir).
			WithWorkdir("/src").
			WithExec([]string{"ls", "-ltrh"}).
			WithEnvVariable("GOOS", "linux").
			WithEnvVariable("GOARCH", "amd64").
			WithEnvVariable("CGO_ENABLED", "0").
			WithExec([]string{"go", "build", "-o",
				"/src/secrets-manager-rotator-lambda"}) // since 'cmd' is where the main.go is.

		if err != nil {
			msg.ShowError("", "Failed to compile the binary", err)
			return err
		}

		outputPath := fmt.Sprintf("%s/%s", platform, config.LambdaName)

		msg.ShowInfo("", fmt.Sprintf("Exporting the binary to: %s", outputPath))

		binaryOutDir = binaryOutDir.WithFile(
			outputPath,
			builder.File(fmt.Sprintf("/src/%s", config.LambdaName)),
		)
	}

	// It exports the binary into a canonical directory: <root>/output/lambda-bin/<platform
	//>/<architecture>/secrets-manager-rotator-lambda
	_, err = binaryOutDir.Export(ctx, config.GetBinaryExportPath())
	if err != nil {
		msg.ShowError("", "Failed to export the binary", err)
		return err
	}

	msg.ShowSuccess("", fmt.Sprintf("The binary has been exported to: %s", config.GetBinaryExportPath()))

	defer client.Close()

	return nil
}
