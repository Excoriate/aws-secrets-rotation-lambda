package lambda

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/errors"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
)

type Compiler interface {
	Compile(srcDir, outputDIr string) (*dagger.Directory, error)
}

type Compile struct {
	Client *dagger.Client
	Ctx    context.Context
}

func (c *Compile) Compile(srcDir, outputDir string) (*dagger.Directory, error) {
	workDir := c.Client.Host().Directory(srcDir)
	msg := tui.NewTUIMessage()

	// Output dir
	binaryOutDir := c.Client.Directory()

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
	}

	_, err := binaryOutDir.Export(c.Ctx, config.GetBinaryExportPath())
	if err != nil {
		return nil, errors.NewTaskError("Failed to export the binary", err)
	}

	return binaryOutDir, nil
}

func NewCompiler(c *dagger.Client, ctx context.Context) Compiler {
	return &Compile{
		Client: c,
		Ctx:    ctx,
	}
}
