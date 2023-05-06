package lambda

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/erroer"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
	"os"
)

type Compiler interface {
	Compile(srcDir, outputDIr string) (*dagger.Directory, string, error)
	Zip(sourceFile, targetFile, targetDir string) (*os.File, string, error)
}

type Compile struct {
	Client *dagger.Client
	Ctx    context.Context
}

func (c *Compile) Zip(sourceFile, targetFile, targetDir string) (*os.File, string, error) {
	if err := common.FileExist(sourceFile); err != nil {
		return nil, "", erroer.NewTaskError("Source file does not exist", err)
	}

	targetDir = fmt.Sprintf("%s/linux/amd64",
		targetDir) // FIXME: Change this logic when multiple-platforms get supported.

	if err := common.DirExist(targetDir); err != nil {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return nil, "", erroer.NewTaskError("Failed to create target directory", err)
		}
	} else {
		if err := os.RemoveAll(targetDir); err != nil {
			return nil, "", erroer.NewTaskError("Failed to remove target directory", err)
		}
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return nil, "", erroer.NewTaskError("Failed to create target directory", err)
		}
	}

	zipFile, err := common.CreateZipFile(sourceFile, targetFile, targetDir)
	if err != nil {
		return nil, "", erroer.NewTaskError("Failed to create zip file", err)
	}

	return zipFile, fmt.Sprintf("%s/%s", targetDir, targetFile), nil
}

func (c *Compile) Compile(srcDir, outputDir string) (*dagger.Directory, string, error) {
	workDir := c.Client.Host().Directory(srcDir)
	msg := tui.NewTUIMessage()

	// Output dir
	binaryOutDir := c.Client.Directory()
	var binaryExportedPath string

	// Compiling the binary.
	platforms := []dagger.Platform{"linux/amd64"}
	for _, platform := range platforms {

		msg.ShowInfo("", fmt.Sprintf("Compiling the binary for platform: %s", platform))

		image := "golang:1.20"
		builder := c.Client.Container(dagger.ContainerOpts{Platform: platform}).
			From(image).
			WithMountedDirectory("/src", workDir).
			WithWorkdir("/src").
			WithExec([]string{"ls", "-ltrh"}).
			WithEnvVariable("GOOS", "linux").
			WithEnvVariable("GOARCH", "amd64").
			WithEnvVariable("CGO_ENABLED", "0").
			WithExec([]string{"go", "build", "-o", fmt.Sprintf("/src/%s", config.LambdaName)})

		outputPath := fmt.Sprintf("%s/%s", platform, config.LambdaName)

		msg.ShowInfo("", fmt.Sprintf("Exporting the binary to: %s", outputPath))

		binaryOutDir = binaryOutDir.WithFile(
			outputPath,
			builder.File(fmt.Sprintf("/src/%s", config.LambdaName)),
		)

		binaryExportedPath = outputPath
	}

	_, err := binaryOutDir.Export(c.Ctx, config.GetBinaryExportPath())
	if err != nil {
		return nil, "", erroer.NewTaskError("Failed to export the binary", err)
	}

	return binaryOutDir, binaryExportedPath, nil
}

func NewCompiler(c *dagger.Client, ctx context.Context) Compiler {
	return &Compile{
		Client: c,
		Ctx:    ctx,
	}
}
