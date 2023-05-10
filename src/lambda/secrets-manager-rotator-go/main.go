package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/rotation"
	"go.uber.org/zap"
	"os"
)

// GetLogger returns a new zap logger
func GetLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

// ListSecretsInAccount lists all secrets in the account
func ListSecretsInAccount(rotator *rotation.RotatorClient) ([]rotation.DiscoveredSecrets,
	error) {

	rotator.Logger.Info("Discovering secrets in account...")
	secrets, err := rotator.Client.ListAll()

	if err != nil {
		rotator.Logger.Error("failed to list secrets", zap.Error(err))
		return nil, fmt.Errorf("failed to list secrets: %v", err)
	}

	var discoveredSecrets []rotation.DiscoveredSecrets

	for _, secret := range secrets {
		secretName := *secret.Name
		secretARN := *secret.ARN
		var rotationIsEnabled bool

		// Log the secret listed
		rotator.Logger.Info(fmt.Sprintf("Discovered secretName: %s", secretName))
		rotator.Logger.Info(fmt.Sprintf("Discovered SecretARN: %s", secretARN))

		// Check if rotation is enabled
		if secret.RotationEnabled != nil {
			rotationIsEnabled = *secret.RotationEnabled
			if rotationIsEnabled {
				rotator.Logger.Info("RotationEnabled: true")
			} else {
				rotator.Logger.Info("RotationEnabled (set but disabled): false")
			}
		} else {
			rotationIsEnabled = false
			rotator.Logger.Info("RotationEnabled: false")
		}

		var secretDescription string
		if secret.Description != nil {
			secretDescription = *secret.Description
		}

		discoveredSecrets = append(discoveredSecrets, rotation.DiscoveredSecrets{
			SecretName:        secretName,
			SecretARN:         secretARN,
			IsRotationEnabled: rotationIsEnabled,
			SecretDescription: secretDescription,
		})
	}

	rotator.Logger.Info("Secrets discovered", zap.Int("count", len(secrets)))
	return discoveredSecrets, nil
}

func handleRequest(ctx context.Context, event rotation.Event) (string, error) {
	logger := GetLogger()
	defer logger.Sync()

	if event == (rotation.Event{}) {
		logger.Fatal("No event received")
	}

	// Logging the event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	logger.Info("Event received for a secret rotation attempt: ", zap.String("event",
		string(eventJson)))

	// Feature flag, enable/disable it setting the environment variable TF_VAR_rotation_lambda_enabled
	isEnabled := os.Getenv("TF_VAR_rotation_lambda_enabled")
	if isEnabled == "false" {
		logger.Fatal("Rotation lambda is disabled")
	}

	// Create the rotator client.
	c, err := rotation.NewRotator(event, logger)

	if err != nil {
		logger.Fatal("AWS Secrets manager rotator lambda cannot be initialised", zap.Error(err))
	}

	// Run pre-checks for rotating this secret
	if err := c.IsRotationAttemptValid(event); err != nil {
		logger.Fatal("Rotation attempt is not valid", zap.Error(err))
	}

	secretId := *event.Arn
	rotationStep := *event.Step
	token := *event.Token
	logger.Info(fmt.Sprintf("Rotation attempt is valid for secret id %s", secretId))
	logger.Info(fmt.Sprintf("Rotation attempt is valid for step %s", rotationStep))
	logger.Info(fmt.Sprintf("Rotation attempt is valid for token %s", token))

	targetSecret, valErr := c.IsSecretValidToRotate(secretId, token)
	if valErr != nil {
		logger.Fatal("Secret is not valid to rotate", zap.Error(valErr))
	}

	// Perform rotation.
	// TODO: Perform an opinionated approach to identify the secret type
	if err := c.Rotate(event, targetSecret, rotationStep, "static"); err != nil {
		logger.Fatal("Secret rotation failed", zap.Error(err))
	}

	return "Secret rotation completed", nil
}

func main() {
	lambda.Start(handleRequest)
}
