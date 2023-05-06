package daggerio

import (
	"dagger.io/dagger"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
)

func SetEnvVars(c *dagger.Container, envVars common.EnvVars) *dagger.Container {
	for k, v := range envVars {
		c = c.WithEnvVariable(k, v)
	}

	return c
}

func ScanAWSCredsFromEnv() common.EnvVars {
	keys := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	}
	return common.GetEnvVars(keys)
}
