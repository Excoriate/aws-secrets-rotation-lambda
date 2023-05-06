package lambda

import (
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/errors"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
)

func IsLambdaSRCDirValid(srcPath string, dirs config.DirConfig) (string, error) {
	msg := tui.NewTUIMessage()

	if !common.IsNotNilAndNotEmpty(srcPath) {
		msg.ShowWarning("", fmt.Sprintf("The lambda source code directory is not set. Using the current directory: %s", dirs.CurrentDir))
		srcPath = dirs.CurrentDir
	}

	if err := common.DirIsValid(srcPath); err != nil {
		return "", errors.NewPipelineConfigurationError(fmt.Sprintf(
			"The lambda source code directory is not a valid directory"+
				": %s. Current directory: %s", srcPath, dirs.CurrentDir), err)
	}

	if err := common.DirIsNotEmpty(srcPath); err != nil {
		return "", errors.NewPipelineConfigurationError(fmt.Sprintf("The lambda source code directory is empty: %s", srcPath), err)
	}

	absDirPath, _ := common.GetDirAbsolute(srcPath)
	srcPath = absDirPath

	return srcPath, nil
}
