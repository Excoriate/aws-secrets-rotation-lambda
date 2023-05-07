package tasks

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/erroer"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
)

func SetInfraConfigInContainer(c *dagger.Container) (*dagger.Container, error) {
	cfg := config.Cfg{}

	envVarKeys := []string{
		"TF_STATE_BUCKET_REGION",
		"TF_STATE_BUCKET",
		"TF_STATE_LOCK_TABLE",
		"TF_REGISTRY_GITHUB_ORG",
		"TF_REGISTRY_BASE_URL",
		"TF_VAR_aws_region",
		"TF_VAR_environment",
		"TF_VAR_rotator_lambda_name",
		"TF_VAR_secret_name",
		"TF_VERSION",
		"TG_VERSION",
	}

	infraEnvVars := make(map[string]string)

	for _, key := range envVarKeys {
		valueFromCLI, err := cfg.GetFromViper(key)

		if err != nil {
			valueFromEnv, err := cfg.GetFromEnvVars(key)
			if err != nil {
				return nil, fmt.Errorf("failed to get value for %s: %v", key, err)
			}
			infraEnvVars[key] = common.RemoveDoubleQuotes(valueFromEnv.Value.(string))
		} else {
			infraEnvVars[key] = common.RemoveDoubleQuotes(valueFromCLI.Value.(string))
		}
	}

	// Add the infraEnvVars to the dagger.Container
	awsEnvVars := daggerio.ScanAWSCredsFromEnv()
	envVarsToSet := common.MergeEnvVars(infraEnvVars, awsEnvVars)

	c = daggerio.SetEnvVars(c, envVarsToSet)

	return c, nil
}

func IsComponentAllowed() (string, error) {
	components := []string{"lambda", "bucket", "secret", "data"}
	cfg := config.Cfg{}
	componentCfg, _ := cfg.GetFromViper("component")

	for _, c := range components {
		if c == componentCfg.Value.(string) {
			return componentCfg.Value.(string), nil
		}
	}

	return "", erroer.NewPipelineConfigurationError("Component is not allowed, "+
		"must be of type: lambda, bucket, secret or data", nil)
}

func IsAllOptionEnabled() bool {
	cfg := config.Cfg{}
	allCfg, _ := cfg.GetFromViperBool("all")

	if allCfg {
		return true
	}

	return false
}

func InfraPlan() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ctx := context.Background()

	tgContainer, component, err := SetupInfra(ctx)
	if err != nil {
		return err
	}

	ux.ShowSubTitle("infra:", fmt.Sprintf("Plan-%s", common.NormaliseStringLower(component)))

	var cmdToRun []string

	if IsAllOptionEnabled() {
		cmdToRun = []string{"terragrunt", "run-all", "plan"}
	} else {
		cmdToRun = []string{"terragrunt", "plan"}

	}
	_, err = tgContainer.
		WithExec([]string{"ls", "-ltrh"}).
		WithExec(cmdToRun).ExitCode(ctx)

	if err != nil {
		msg.ShowError("", "Terragrunt plan failed", err)
		return err
	}

	msg.ShowSuccess("", "Terragrunt plan succeeded")
	return nil
}

func InfraDeploy() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ctx := context.Background()

	tgContainer, component, err := SetupInfra(ctx)
	if err != nil {
		return err
	}

	ux.ShowSubTitle("infra:", fmt.Sprintf("Apply-%s", common.NormaliseStringLower(component)))

	var cmdToRun []string
	if IsAllOptionEnabled() {
		cmdToRun = []string{"terragrunt", "run-all", "apply", "-auto-approve"}
	} else {
		cmdToRun = []string{"terragrunt", "apply", "-auto-approve"}
	}
	_, err = tgContainer.
		WithExec([]string{"ls", "-ltrh"}).
		WithExec(cmdToRun).ExitCode(ctx)

	if err != nil {
		msg.ShowError("", "Terragrunt apply failed", err)
		return err
	}

	msg.ShowSuccess("", "Terragrunt apply succeeded")
	return nil
}

func InfraDestroy() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ctx := context.Background()

	tgContainer, component, err := SetupInfra(ctx)
	if err != nil {
		return err
	}

	ux.ShowSubTitle("infra:", fmt.Sprintf("Destroy-%s", common.NormaliseStringLower(component)))

	var cmdToRun []string
	if IsAllOptionEnabled() {
		cmdToRun = []string{"terragrunt", "run-all", "destroy", "-auto-approve"}
	} else {
		cmdToRun = []string{"terragrunt", "destroy", "-auto-approve"}
	}
	_, err = tgContainer.
		WithExec([]string{"ls", "-ltrh"}).
		WithExec(cmdToRun).ExitCode(ctx)

	if err != nil {
		msg.ShowError("", "Terragrunt destroy failed", err)
		return err
	}

	msg.ShowSuccess("", "Terragrunt destroy succeeded")
	return nil
}

func SetupInfra(ctx context.Context) (*dagger.Container, string, error) {
	msg := tui.NewTUIMessage()

	// Validating component
	component, err := IsComponentAllowed()
	if err != nil {
		msg.ShowError("", "Failed to validate component", err)
		return nil, "", err
	}

	// Getting directories configuration
	dirs, dirsErr := config.GetDirConfig()
	if dirsErr != nil {
		msg.ShowError("", "Failed to get directories configuration", err)
		return nil, "", dirsErr
	}

	// Booting dagger!
	client, engineErr := daggerio.NewClientWithWorkDir(ctx,
		dirs.GitRepoDir) // It's required to set the git root dir, since TG depends on discovery of the root dir.
	if engineErr != nil {
		msg.ShowError("", "Failed to boot dagger", err)
		return nil, "", engineErr
	}

	// Bootstrap container
	image := "alpine/terragrunt"

	toMount := client.Host().Directory(".")

	tgContainer := client.Container().
		From(image).
		WithMountedDirectory("/src", toMount).
		WithWorkdir(fmt.Sprintf("/src/%s", config.GetInfraWorkDirPath(component)))

	// Add required environment variables.
	tgContainer, err = SetInfraConfigInContainer(tgContainer)
	if err != nil {
		msg.ShowError("", "Failed to set environment variables in container", err)
		return nil, "", err
	}

	return tgContainer, component, nil
}
