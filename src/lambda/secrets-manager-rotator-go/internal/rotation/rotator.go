package rotation

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/adapter"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/client"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/erroer"
	"go.uber.org/zap"
)

type Rotator interface {
	IsRotationAttemptValid(event Event) error
	IsSecretValidToRotate(secretArn, token string) (*secretsmanager.DescribeSecretOutput, error)
	Rotate(event Event, secret *secretsmanager.DescribeSecretOutput, step string,
		secretType string) error
}

type RotatorClient struct {
	Logger         *zap.Logger
	Client         client.SecretsManager
	SecretToRotate Event
}

func (r *RotatorClient) Rotate(event Event, secret *secretsmanager.DescribeSecretOutput, step string,
	secretType string) error {

	steps := GetSteps()

	secretId := *secret.ARN
	secretName := *secret.Name
	secretToken := *event.Token

	r.Logger.Info(fmt.Sprintf("Initializing rotation for secret: %s, "+
		"with id: %s on step: %s with token: "+secretName, secretId, step, secretToken))

	s := NewStepExecutionerClient(r.Logger, r.Client, event, secret)

	switch step {
	case steps.Create:
		return s.CreateSecretStep()
	case steps.Set:
		return s.SetSecretStep()
	case steps.Test:
		return s.TestSecretStep()
	case steps.Finish:
		return s.FinishSecretStep()
	}

	return nil
}

func (r *RotatorClient) IsRotationAttemptValid(event Event) error {
	if event == (Event{}) {
		r.Logger.Error("Event is empty")
		return erroer.NewValidationError("event is empty", nil)
	}

	secretArn := *event.Arn

	if secretArn == "" {
		r.Logger.Error("SecretArn is empty")
		return erroer.NewValidationError("secretArn is empty", nil)
	}

	token := *event.Token

	if token == "" {
		r.Logger.Error("ClientRequestToken is empty")
		return erroer.NewValidationError("clientRequestToken is empty", nil)
	}

	step := *event.Step
	if step == "" {
		r.Logger.Error("Step is empty")
		return erroer.NewValidationError("step is empty", nil)
	}

	// If the step isn't found, fail
	isFound := false
	for _, allowedStep := range AllowedSteps {
		if step == allowedStep {
			isFound = true
		}
	}

	if !isFound {
		r.Logger.Error(fmt.Sprintf("Step is not valid: %s", step))
		return erroer.NewValidationError(fmt.Sprintf("step is not valid: %s", step), nil)
	}

	return nil
}

func (r *RotatorClient) IsSecretValidToRotate(secretId, token string) (*secretsmanager.
DescribeSecretOutput, error) {
	secret, err := r.Client.GetSecret(secretId)

	if err != nil {
		r.Logger.Error(fmt.Sprintf("Error getting secret with arn: %s", secretId), zap.Error(err))
		return nil, erroer.NewValidationError(fmt.Sprintf("error getting secret with arn: %s",
			secretId), err)
	}

	if secret.RotationEnabled != nil {
		rotationIsEnabled := *secret.RotationEnabled
		if rotationIsEnabled {
			r.Logger.Info("RotationEnabled: true")
		} else {
			r.Logger.Error("RotationEnabled (set but disabled): false")
			return nil, erroer.NewValidationError(fmt.Sprintf(
				"rotation is not enabled for secret: %s", secretId), nil)
		}
	} else {
		r.Logger.Error("RotationEnabled: false")
		return nil, erroer.NewValidationError(fmt.Sprintf("rotation is not enabled for secret: %s",
			secretId), nil)
	}

	// Check labels
	if _, ok := secret.VersionIdsToStages[token]; !ok {
		r.Logger.Error(fmt.Sprintf("Secret version %s has no stage for rotation of secret %s.",
			token, secretId))
		return nil, erroer.NewSecretError(fmt.Sprintf("secret version %s has no stage for rotation of secret"+
			" %s.",
			token, secretId), nil)
	}

	stagingLabels := GetStagingLabels()
	for _, value := range secret.VersionIdsToStages[token] {
		// If the secret is already the 'AWS_CURRENT' version, fail
		if value == stagingLabels.Current {
			r.Logger.Error(fmt.Sprintf("Secret version %s already set as %s for secret %s.", token,
				stagingLabels.Current, secretId))
			return nil, erroer.NewSecretError(fmt.Sprintf(
				"secret version %s already set as %s for secret"+" %s.", token,
				stagingLabels.Current, secretId), nil)
		}

		// If the secret isn't set with the 'AWSPENDING' label, fail.
		if value != stagingLabels.Pending {
			r.Logger.Error(fmt.Sprintf("Secret version %s not set as %s for secret %s.", token,
				stagingLabels.Pending, secretId))
			return nil, erroer.NewSecretError(fmt.Sprintf(
				"secret version %s not set as %s for secret"+" %s.", token,
				stagingLabels.Pending, secretId), nil)
		}
	}

	// Checking if the secret version with the AWSCURRENT stage label exists is a validation step
	// to ensure that there's a current version of the secret in use before proceeding with the
	// rotation process. This helps to avoid potential issues where the rotation process could be
	// accidentally triggered on a secret that has not been properly set up or has no existing
	// version marked as current.
	//
	// By confirming the presence of a secret version with the AWSCURRENT stage label,
	// the rotation process can proceed safely and update the secret,
	// knowing that there's a working version of the secret that the application or service is
	// already using. This check also helps to maintain the integrity of the secret rotation
	//process and avoid unexpected issues related to missing or improperly configured secrets.
	currentSecretVersion, currentVersionErr := r.Client.GetSecretValueByStageLabel(secretId, "",
		stagingLabels.Current)
	if currentVersionErr != nil {
		r.Logger.Error(fmt.Sprintf("This secret %s can not be rotated because there is no version"+
			" "+"present with AWSCURRENT stage label", secretId))

		return nil, erroer.NewSecretError(fmt.Sprintf(
			"this secret %s can not be rotated because there is no version"+" "+"present with"+
				" AWSCURRENT stage label", secretId), currentVersionErr)
	}

	r.Logger.Info(fmt.Sprintf("Successfully found a version with AWSCURRENT stage label in secret %s, with token (versionId) %s", secretId, currentSecretVersion.VersionId))
	return secret, nil
}

func NewRotator(event Event, logger *zap.Logger) (*RotatorClient, error) {
	// Get adapter client
	awsCfg, err := adapter.NewAWS("")

	if err != nil {
		logger.Error("Failed to initialise Rotator Client. Can't instantiate AWS client", zap.Error(err))
		return nil, erroer.NewConfigurationError("failed to initialise Rotator Client. Can't instantiate AWS client: %v", err)
	}

	// Get Secrets Manager client
	smClient := client.NewSecretsManager(awsCfg, logger)

	logger.Info("Rotator client initialised")

	return &RotatorClient{
		Logger:         logger,
		Client:         smClient,
		SecretToRotate: event,
	}, nil

}
