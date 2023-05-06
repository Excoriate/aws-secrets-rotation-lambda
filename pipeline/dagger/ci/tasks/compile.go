package tasks

import (
	"context"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/lambda"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
)

func CompileLambda() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ux.ShowSubTitle("lambda:", "Compile")

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

	// Validating the Lambda source code directory.
	//if err != nil {
	//	msg.ShowWarning("", fmt.Sprintf("The lambda source code directory is not set. Using the current directory: %s", dirs.CurrentDir))
	//	srcPath = dirs.CurrentDir
	//} else {
	//	srcPath = cfgValue.Value.(string)
	//
	//	if err := common.DirIsValid(srcPath); err != nil {
	//		msg.ShowError("",
	//			fmt.Sprintf("The lambda source code directory is not a valid directory"+
	//				": %s. Current directory: %s", srcPath, dirs.CurrentDir), err)
	//
	//		return err
	//	}
	//
	//	if err := common.DirIsNotEmpty(srcPath); err != nil {
	//		msg.ShowError("", fmt.Sprintf("The lambda source code directory is empty: %s", srcPath), err)
	//		return err
	//	}
	//
	//	absDirPath, _ := common.GetDirAbsolute(srcPath)
	//	srcPath = absDirPath
	//
	//	msg.ShowInfo("", fmt.Sprintf("Using the lambda source code directory: %s", srcPath))
	//}

	// Compiling
	_, err = compiler.Compile(srcPath, "")
	if err != nil {
		msg.ShowError("", "The binary could not be compiled", err)
		return err
	}

	msg.ShowSuccess("", fmt.Sprintf("The binary has been exported to: %s", config.GetBinaryExportPath()))

	defer client.Close()

	return nil
}
